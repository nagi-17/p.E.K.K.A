package services

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/nagi-17/p.E.K.K.A/internal/database"
	"github.com/nagi-17/p.E.K.K.A/internal/models"
)

func PlaceNewBuilding(ctx context.Context, playerID uuid.UUID, bType string, x int, y int) error {
	var bData *models.BuildingData
	bData, err := models.GetBuildingDataByTypeLevel(ctx, bType, 1)
	if err != nil {
		return ErrInvalidBuildingType
	}

	townHallLevel, err := models.GetPlayerTownHallLevel(ctx, playerID)
	if err != nil {
		return err
	}
	switch bType {
	case "Cannon", "Archer Tower", "Mortar":
		defData, err := models.GetDefBuildingData(ctx, bType, bData.BuildingLevel)
		if err != nil {
			return err
		}
		if townHallLevel < defData.UnlockTownHallLevel {
			return ErrTownHallUnderLevel
		}
	case "Elixir Collector", "Pancake Machine":
		resData, err := models.GetResBuildingData(ctx, bType, bData.BuildingLevel)
		if err != nil {
			return err
		}
		if townHallLevel < resData.UnlockTownHallLevel {
			return ErrTownHallUnderLevel
		}
	case "Elixir Storage", "Pancake Stack":
		strgData, err := models.GetStrgBuildingData(ctx, bType, bData.BuildingLevel)
		if err != nil {
			return err
		}
		if townHallLevel < strgData.UnlockTownHallLevel {
			return ErrTownHallUnderLevel
		}
	case "Laboratory":
		labData, err := models.GetLabData(ctx, bType, bData.BuildingLevel)
		if err != nil {
			return err
		}
		if townHallLevel < labData.UnlockTownHallLevel {
			return ErrTownHallUnderLevel
		}
	case "Army Camp":
		campData, err := models.GetArmyCampData(ctx, bType, bData.BuildingLevel)
		if err != nil {
			return err
		}
		if townHallLevel < campData.UnlockTownHallLevel {
			return ErrTownHallUnderLevel
		}
	default:
		return ErrInvalidBuildingType
	}

	tx, err := database.DB.Begin(ctx)
	if err != nil {
		return fmt.Errorf("Failed to begin database transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	valid := models.IsCellValid(ctx, playerID, x, y, bData.Width, bData.Height, uuid.Nil)
	if valid != nil {
		if valid.Error() == "Cell is occupied" {
			return ErrCellOccupied
		}
		return valid
	}

	playerRes, err := models.GetPlayerInfoByID(ctx, playerID)
	if err != nil {
		return fmt.Errorf("Error fetching player resources stats: %w", err)
	}
	elixir := playerRes.Elixir
	pancakes := playerRes.Pancakes
	if pancakes < bData.UpgradeCostPancakes {
		return ErrNotEnoughPancakes
	}

	if elixir < bData.UpgradeCostElixir {
		return ErrNotEnoughElixir
	}

	allOwnedBuildings, err := models.GetOwnedBuildingData(ctx, playerID)
	if err != nil {
		return fmt.Errorf("Failed to fetch owned building data: %w", err)
	}
	var count int = 0
	for i := 0; i < len(allOwnedBuildings); i++ {
		if allOwnedBuildings[i].BuildingType == bType {
			count++
		}
	}
	if (count + 1) > bData.MaxQuantityAvailable {
		return ErrBuildingLimitReached
	}

	var upgrade_complete_at *time.Time = nil
	var last_collected_at *time.Time = nil

	if bData.BuildTime > 0 {
		finishTime := time.Now().Add(time.Duration(bData.BuildTime) * time.Second)
		upgrade_complete_at = &finishTime
		sentinel := time.Unix(0, 0)
		last_collected_at = &sentinel
	} else {
		if bType == "Elixir Collector" || bType == "Pancake Machine" {
			currTime := time.Now()
			last_collected_at = &currTime
		}
	}
	err = models.InsertBuilding(ctx, tx, playerID, bData.ID, x, y, upgrade_complete_at, last_collected_at)
	if err != nil {
		return fmt.Errorf("Error in placing new building: %w", err)
	}

	newSkill := playerRes.Skill_points + bData.SkillOnUpgrade
	err = models.AddSkill(ctx, tx, playerID, newSkill)
	if err != nil {
		return err
	}

	newPancakes := pancakes - bData.UpgradeCostPancakes
	newElixir := elixir - bData.UpgradeCostElixir
	err = models.UpdateResources(ctx, tx, playerID, newPancakes, newElixir)
	if err != nil {
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("Failed to commit to database: %w", err)
	}

	return nil
}

func StartUpgrade(ctx context.Context, ownedBuildingID uuid.UUID) error {
	var bDataID int
	var playerID uuid.UUID
	var curr_upgrade_complete_time *time.Time

	bDataID, playerID, curr_upgrade_complete_time, err := models.GetBuildingAndPlayerIDUsingOwnedBuildingID(ctx, ownedBuildingID)
	if err != nil {
		return fmt.Errorf("Cannot fetch building id: %w", err)
	}
	if curr_upgrade_complete_time != nil {
		if time.Now().Before(*curr_upgrade_complete_time) {
			return fmt.Errorf("Building already under upgrade")
		}
	}

	bData, err := models.GetBuildingDataByID(ctx, bDataID)
	if err != nil {
		return fmt.Errorf("Can't fetch building data: %w", err)
	}

	playerStats, err := models.GetPlayerInfoByID(ctx, playerID)
	if err != nil {
		return fmt.Errorf("Error in fetching player stats: %w", err)
	}
	if playerStats.Elixir < bData.UpgradeCostElixir || playerStats.Pancakes < bData.UpgradeCostPancakes {
		return fmt.Errorf("Can't upgrade building: insufficient resources")
	}

	townHallLevel, err := models.GetPlayerTownHallLevel(ctx, playerID)
	if err != nil {
		return err
	}

	switch bData.BuildingType {
	case "Town Hall":
		if bData.BuildingLevel >= 4 {
			return fmt.Errorf("Town Hall is already maxed out")
		}

		x, err := models.GetTownHallData(ctx, bData.BuildingLevel)
		if err != nil {
			return fmt.Errorf("Can't fetch town hall data: %w", err)
		}

		if playerStats.Skill_points < x.MinSkillPointsBeforeUpgrade {
			return fmt.Errorf("Can't upgrade Town Hall: %d skill points required to upgrade", x.MinSkillPointsBeforeUpgrade)
		}

	case "Cannon", "Archer Tower", "Mortar":
		x, err := models.GetDefBuildingData(ctx, bData.BuildingType, bData.BuildingLevel)
		if err != nil {
			return fmt.Errorf("Error in fetching defense building stats: %w", err)
		}

		if bData.BuildingLevel == x.MaxPossibleUpgradeLevel {
			return fmt.Errorf("Building is already maxed out")
		}

		nextDef, err := models.GetDefBuildingData(ctx, bData.BuildingType, bData.BuildingLevel+1)
		if err != nil {
			return fmt.Errorf("Error in fetching next level stats: %w", err)
		}
		if townHallLevel < nextDef.UnlockTownHallLevel {
			return fmt.Errorf("Town Hall level %d required to upgrade this building", nextDef.UnlockTownHallLevel)
		}

	case "Elixir Collector", "Pancake Machine":
		x, err := models.GetResBuildingData(ctx, bData.BuildingType, bData.BuildingLevel)
		if err != nil {
			return fmt.Errorf("Can't fetch resource building stats: %w", err)
		}

		if bData.BuildingLevel == x.MaxPossibleUpgradeLevel {
			return fmt.Errorf("Building is already maxed out")
		}

		nextRes, err := models.GetResBuildingData(ctx, bData.BuildingType, bData.BuildingLevel+1)
		if err != nil {
			return fmt.Errorf("Error in fetching next level stats: %w", err)
		}
		if townHallLevel < nextRes.UnlockTownHallLevel {
			return fmt.Errorf("Town Hall level %d required to upgrade this building", nextRes.UnlockTownHallLevel)
		}

	case "Elixir Storage", "Pancake Stack":
		x, err := models.GetStrgBuildingData(ctx, bData.BuildingType, bData.BuildingLevel)
		if err != nil {
			return fmt.Errorf("Can't fetch storage building stats: %w", err)
		}

		if bData.BuildingLevel == x.MaxPossibleUpgradeLevel {
			return fmt.Errorf("Building is already maxed out")
		}

		nextStrg, err := models.GetStrgBuildingData(ctx, bData.BuildingType, bData.BuildingLevel+1)
		if err != nil {
			return fmt.Errorf("Error in fetching next level stats: %w", err)
		}
		if townHallLevel < nextStrg.UnlockTownHallLevel {
			return fmt.Errorf("Town Hall level %d required to upgrade this building", nextStrg.UnlockTownHallLevel)
		}

	case "Laboratory":
		x, err := models.GetLabData(ctx, bData.BuildingType, bData.BuildingLevel)
		if err != nil {
			return fmt.Errorf("Can't fetch laboratory stats: %w", err)
		}

		if bData.BuildingLevel == x.MaxPossibleUpgradeLevel {
			return fmt.Errorf("Building is already maxed out")
		}

		nextLab, err := models.GetLabData(ctx, bData.BuildingType, bData.BuildingLevel+1)
		if err != nil {
			return fmt.Errorf("Error in fetching next level stats: %w", err)
		}
		if townHallLevel < nextLab.UnlockTownHallLevel {
			return fmt.Errorf("Town Hall level %d required to upgrade this building", nextLab.UnlockTownHallLevel)
		}

	case "Army Camp":
		x, err := models.GetArmyCampData(ctx, bData.BuildingType, bData.BuildingLevel)
		if err != nil {
			return fmt.Errorf("Can't fetch army camp stats: %w", err)
		}

		if bData.BuildingLevel == x.MaxPossibleUpgradeLevel {
			return fmt.Errorf("Building is already maxed out")
		}

		nextCamp, err := models.GetArmyCampData(ctx, bData.BuildingType, bData.BuildingLevel+1)
		if err != nil {
			return fmt.Errorf("Error in fetching next level stats: %w", err)
		}
		if townHallLevel < nextCamp.UnlockTownHallLevel {
			return fmt.Errorf("Town Hall level %d required to upgrade this building", nextCamp.UnlockTownHallLevel)
		}
	}

	if bData.UpgradeTime == 0 {
		upgradedLevel := bData.BuildingLevel + 1
		upgradedBData, err := models.GetBuildingDataByTypeLevel(ctx, bData.BuildingType, upgradedLevel)
		if err != nil {
			return fmt.Errorf("Failed to find data of upgraded building: %w", err)
		}

		tx, err := database.DB.Begin(ctx)
		if err != nil {
			return fmt.Errorf("Failed to begin database transaction: %w", err)
		}
		defer tx.Rollback(ctx)

		queryUpdate := `UPDATE owned_building SET building_data_id = $1, upgrade_complete_at = NULL WHERE id = $2`
		_, err = tx.Exec(ctx, queryUpdate, upgradedBData.ID, ownedBuildingID)
		if err != nil {
			return fmt.Errorf("Error in updating owned building data: %w", err)
		}

		newPancakes := playerStats.Pancakes - bData.UpgradeCostPancakes
		newElixir := playerStats.Elixir - bData.UpgradeCostElixir
		err = models.UpdateResources(ctx, tx, playerID, newPancakes, newElixir)
		if err != nil {
			return fmt.Errorf("Error updating player stats: %w", err)
		}

		newSkill := playerStats.Skill_points + bData.SkillOnUpgrade
		err = models.AddSkill(ctx, tx, playerID, newSkill)
		if err != nil {
			return fmt.Errorf("Error updating player skill: %w", err)
		}

		err = tx.Commit(ctx)
		if err != nil {
			return fmt.Errorf("Failed to commit to database: %w", err)
		}
		return nil
	}

	tx, err := database.DB.Begin(ctx)
	if err != nil {
		return fmt.Errorf("Failed to begin database transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	finishTime := time.Now().Add(time.Duration(bData.UpgradeTime) * time.Second)

	err = models.StartBuildingUpgradeQueryFunc(ctx, tx, ownedBuildingID, finishTime)
	if err != nil {
		return fmt.Errorf("Error in updating finish time of building: %w", err)
	}

	newPancakes := playerStats.Pancakes - bData.UpgradeCostPancakes
	newElixir := playerStats.Elixir - bData.UpgradeCostElixir
	err = models.UpdateResources(ctx, tx, playerID, newPancakes, newElixir)
	if err != nil {
		return fmt.Errorf("Error updating player stats: %w", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("Failed to commit to database: %w", err)
	}

	return nil
}
