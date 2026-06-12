package models

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/nagi-17/p.E.K.K.A/internal/database"
)

func GetOwnedBuildingData(ctx context.Context, playerID uuid.UUID) ([]OwnedBuildingWithData, error) {
	query := `
	SELECT 
		ob.id, ob.player_id, ob.building_data_id, ob.pos_x, ob.pos_y, 
		ob.upgrade_complete_at, ob.last_collected_at,
		bd.building_type, bd.building_level, bd.health, bd.width, bd.height, 
		bd.build_time, bd.upgrade_cost_elixir, bd.upgrade_cost_pancakes, 
		bd.upgrade_time, bd.max_quantity_available, bd.skill_on_upgrade
	FROM owned_building ob
	INNER JOIN building_data bd ON ob.building_data_id = bd.id
	WHERE ob.player_id = $1;
	`
	rows, err := database.DB.Query(ctx, query, playerID)
	if err != nil {
		return nil, fmt.Errorf("Error in fetching owned buildings data: %w", err)
	}
	defer rows.Close()

	var buildings []OwnedBuildingWithData
	for rows.Next() {
		var x OwnedBuildingWithData
		err := rows.Scan(
			&x.OwnedBuildingData.ID, &x.OwnedBuildingData.PlayerID, &x.OwnedBuildingData.BuildingDataID,
			&x.OwnedBuildingData.PosX, &x.OwnedBuildingData.PosY, &x.OwnedBuildingData.UpgradeCompleteAt,
			&x.OwnedBuildingData.LastCollectedAt, &x.BuildingType, &x.BuildingLevel, &x.Health, &x.Width,
			&x.Height, &x.BuildTime, &x.UpgradeCostElixir, &x.UpgradeCostPancakes, &x.UpgradeTime,
			&x.MaxQuantityAvailable, &x.SkillOnUpgrade,
		)
		if err != nil {
			return nil, fmt.Errorf("Error in scanning owned building data/rows: %w", err)
		}
		buildings = append(buildings, x)
	}
	return buildings, nil
}

func PlaceNewBuilding(ctx context.Context, playerID uuid.UUID, bType string, x int, y int) error {
	var bData *BuildingData
	bData, err := GetBuildingDataByTypeLevel(ctx, bType, 1)
	if err != nil {
		return fmt.Errorf("Invalid building type")
	}

	townHallLevel, err := GetPlayerTownHallLevel(ctx, playerID)
	if err != nil {
		return err
	}
	switch bType {
	case "Cannon", "Archer Tower", "Mortar":
		defData, err := GetDefBuildingData(ctx, bType, bData.BuildingLevel)
		if err != nil {
			return err
		}
		if townHallLevel < defData.UnlockTownHallLevel {
			return fmt.Errorf("Town Hall is under levelled")
		}
	case "Elixir Collector", "Pancake Machine":
		resData, err := GetResBuildingData(ctx, bType, bData.BuildingLevel)
		if err != nil {
			return err
		}
		if townHallLevel < resData.UnlockTownHallLevel {
			return fmt.Errorf("Town Hall is under levelled")
		}
	case "Elixir Storage", "Pancake Stack":
		strgData, err := GetStrgBuildingData(ctx, bType, bData.BuildingLevel)
		if err != nil {
			return err
		}
		if townHallLevel < strgData.UnlockTownHallLevel {
			return fmt.Errorf("Town Hall is under levelled")
		}
	case "Laboratory":
		labData, err := GetLabData(ctx, bType, bData.BuildingLevel)
		if err != nil {
			return err
		}
		if townHallLevel < labData.UnlockTownHallLevel {
			return fmt.Errorf("Town Hall is under levelled")
		}
	default:
		return fmt.Errorf("Invalid building type")
	}

	tx, err := database.DB.Begin(ctx)
	if err != nil {
		return fmt.Errorf("Failed to begin databse transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	valid := IsCellValid(ctx, playerID, x, y, bData.Width, bData.Height, uuid.Nil)
	if valid != nil {
		return valid
	}

	playerRes, err := GetPlayerInfoByID(ctx, playerID)
	if err != nil {
		return fmt.Errorf("Error fetching player resources stats: %w", err)
	}
	elixir := playerRes.Elixir
	pancakes := playerRes.Pancakes
	if pancakes < bData.UpgradeCostPancakes {
		return fmt.Errorf("Not enough pancakes")
	}

	if elixir < bData.UpgradeCostElixir {
		return fmt.Errorf("Not enough elixir")
	}

	allOwnedBuildings, err := GetOwnedBuildingData(ctx, playerID)
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
		return fmt.Errorf("All possible buildings of this type have already been placed")
	}

	query := `
	INSERT INTO owned_building (id, player_id, building_data_id, pos_x, pos_y, upgrade_complete_at, last_collected_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	var id uuid.UUID = uuid.New()
	var upgrade_complete_at *time.Time = nil
	var last_collected_at *time.Time = nil
	if bData.BuildTime > 0 {
		finishTime := time.Now().Add(bData.BuildTime)
		upgrade_complete_at = &finishTime
	}
	_, err = tx.Exec(ctx, query, id, playerID, bData.ID, x, y, upgrade_complete_at, last_collected_at)
	if err != nil {
		return fmt.Errorf("Error in placing new building: %w", err)
	}

	newSkill := playerRes.Skill_points + bData.SkillOnUpgrade
	err = AddSkill(ctx, tx, playerID, newSkill)
	if err != nil {
		return err
	}

	newPancakes := pancakes - bData.UpgradeCostPancakes
	newElixir := elixir - bData.UpgradeCostElixir
	err = UpdateResources(ctx, tx, playerID, newPancakes, newElixir)
	if err != nil {
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("Failed to commit to database: %w", err)
	}

	return nil
}

func MoveBuilding(ctx context.Context, ownedBuildingID uuid.UUID, newX int, newY int) error {

	query1 := `
	SELECT ob.player_id, bd.width, bd.height FROM owned_building ob
	INNER JOIN building_data bd ON ob.building_data_id=bd.id
	WHERE ob.id=$1
	`
	var playerID uuid.UUID
	var width, height int
	err := database.DB.QueryRow(ctx, query1, ownedBuildingID).Scan(&playerID, &width, &height)
	if err != nil {
		return fmt.Errorf("Couldn't fetch width, height from building_data table")
	}

	err = IsCellValid(ctx, playerID, newX, newY, width, height, ownedBuildingID)
	if err != nil {
		return err
	}

	tx, err := database.DB.Begin(ctx)
	if err != nil {
		return fmt.Errorf("Failed to begin databse transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	query2 := `
	UPDATE owned_building
	SET pos_x=$1, pos_y=$2
	WHERE owned_building.id=$3
	`
	_, err = tx.Exec(ctx, query2, newX, newY, ownedBuildingID)
	if err != nil {
		return fmt.Errorf("Error in updating new coordinates of building: %w", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("Failed to commit to database: %w", err)
	}

	return nil
}

func IsCellValid(ctx context.Context, playerID uuid.UUID, x int, y int, width int, height int, ignoreBuilding uuid.UUID) error {

	if x < 0 || y < 0 || x+width > 40 || y+height > 40 {
		return fmt.Errorf("Building is out of bounds")
	}

	allOwnedBuildings, err := GetOwnedBuildingData(ctx, playerID)
	if err != nil {
		return fmt.Errorf("Failed to fetch owned building data: %w", err)
	}
	for i := 0; i < len(allOwnedBuildings); i++ {
		if allOwnedBuildings[i].ID == ignoreBuilding {
			continue
		}
		//Reference pt: Bottom left corner
		check_x := x < (allOwnedBuildings[i].PosX+allOwnedBuildings[i].Width) && (x+width) > allOwnedBuildings[i].PosX
		check_y := y < (allOwnedBuildings[i].PosY+allOwnedBuildings[i].Height) && (y+height) > allOwnedBuildings[i].PosY
		if check_x && check_y {
			return fmt.Errorf("Cell is occupied")
		}
	}
	return nil
}

func StartUpgrade(ctx context.Context, ownedBuildingID uuid.UUID) error {
	query1 := `SELECT building_data_id, player_id, upgrade_complete_at FROM owned_building WHERE owned_building.id = $1 `

	var bDataID int
	var playerID uuid.UUID
	var curr_upgrade_complete_time *time.Time

	err := database.DB.QueryRow(ctx, query1, ownedBuildingID).Scan(&bDataID, &playerID, &curr_upgrade_complete_time)
	if err != nil {
		return fmt.Errorf("Cannot fetch building id: %w", err)
	}
	if curr_upgrade_complete_time != nil {
		if time.Now().Before(*curr_upgrade_complete_time) {
			return fmt.Errorf("Building already under upgrade")
		}
	}

	bData, err := GetBuildingDataByID(ctx, bDataID)
	if err != nil {
		return fmt.Errorf("Can't fetch building data: %w", err)
	}

	playerStats, err := GetPlayerInfoByID(ctx, playerID)
	if err != nil {
		return fmt.Errorf("Error in fetching player stats: %w", err)
	}
	if playerStats.Elixir < bData.UpgradeCostElixir || playerStats.Pancakes < bData.UpgradeCostPancakes {
		return fmt.Errorf("Can't upgrade building: insufficient resources")
	}

	switch bData.BuildingType {
	case "Town Hall":
		x, err := GetTownHallData(ctx, bData.BuildingLevel)
		if err != nil {
			return fmt.Errorf("Can't fetch town hall data: %w", err)
		}

		if playerStats.Skill_points < x.MinSkillPointsBeforeUpgrade {
			return fmt.Errorf("Can't upgrade Town Hall: %d skill points required to upgrade", x.MinSkillPointsBeforeUpgrade)
		}

	case "Cannon", "Archer Tower", "Mortar":
		x, err := GetDefBuildingData(ctx, bData.BuildingType, bData.BuildingLevel)
		if err != nil {
			return fmt.Errorf("Error in fetching defense building stats: %w", err)
		}

		if bData.BuildingLevel == x.MaxPossibleUpgradeLevel {
			return fmt.Errorf("Building is already maxed out")
		}

	case "Elixir Collector", "Pancake Machine":
		x, err := GetResBuildingData(ctx, bData.BuildingType, bData.BuildingLevel)
		if err != nil {
			return fmt.Errorf("Can't fetch resource building stats: %w", err)
		}

		if bData.BuildingLevel == x.MaxPossibleUpgradeLevel {
			return fmt.Errorf("Building is already maxed out")
		}

	case "Elixir Storage", "Pancake Stack":
		x, err := GetStrgBuildingData(ctx, bData.BuildingType, bData.BuildingLevel)
		if err != nil {
			return fmt.Errorf("Can't fetch storage building stats: %w", err)
		}

		if bData.BuildingLevel == x.MaxPossibleUpgradeLevel {
			return fmt.Errorf("Building is already maxed out")
		}

	case "Laboratory":
		x, err := GetLabData(ctx, bData.BuildingType, bData.BuildingLevel)
		if err != nil {
			return fmt.Errorf("Can't fetch laboratory stats: %w", err)
		}

		if bData.BuildingLevel == x.MaxPossibleUpgradeLevel {
			return fmt.Errorf("Building is already maxed out")
		}
	}

	tx, err := database.DB.Begin(ctx)
	if err != nil {
		return fmt.Errorf("Failed to begin databse transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	query2 := `UPDATE owned_building SET upgrade_complete_at = $1 WHERE owned_building.id = $2`
	finishTime := time.Now().Add(bData.UpgradeTime)

	_, err = tx.Exec(ctx, query2, finishTime, ownedBuildingID)
	if err != nil {
		return fmt.Errorf("Error in updating finish time of building: %w", err)
	}

	newPancakes := playerStats.Pancakes - bData.UpgradeCostPancakes
	newElixir := playerStats.Elixir - bData.UpgradeCostElixir
	err = UpdateResources(ctx, tx, playerID, newPancakes, newElixir)
	if err != nil {
		return fmt.Errorf("Error updating player stats: %w", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("Failed to commit to database: %w", err)
	}

	return nil
}

func FinishUpgrade(ctx context.Context, ownedBuildingID uuid.UUID) error {
	query1 := `
	SELECT ob.player_id, ob.upgrade_complete_at, bd.building_type, bd.building_level, bd.skill_on_upgrade
	FROM owned_building ob
	INNER JOIN building_data bd ON ob.building_data_id = bd.id
	WHERE ob.id = $1
	`
	var playerID uuid.UUID
	var upgradeCompleteAt *time.Time
	var bType string
	var currLevel int
	var skillGain int
	err := database.DB.QueryRow(ctx, query1, ownedBuildingID).Scan(&playerID, &upgradeCompleteAt, &bType, &currLevel, &skillGain)
	if err != nil {
		return fmt.Errorf("Failed to fetch building: %w", err)
	}

	if upgradeCompleteAt == nil {
		return fmt.Errorf("Building is not under upgrade")
	}
	if time.Now().Before(*upgradeCompleteAt) {
		return fmt.Errorf("Upgrade is still going on")
	}

	upgradedLevel := currLevel + 1
	upgradedBData, err := GetBuildingDataByTypeLevel(ctx, bType, upgradedLevel)
	if err != nil {
		return fmt.Errorf("Failed to find data of upgraded building: %w", err)
	}

	tx, err := database.DB.Begin(ctx)
	if err != nil {
		return fmt.Errorf("Failed to begin databse transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	query2 := `UPDATE owned_building SET building_data_id = $1, upgrade_complete_at = NULL WHERE id = $2`
	_, err = tx.Exec(ctx, query2, upgradedBData.ID, ownedBuildingID)
	if err != nil {
		return fmt.Errorf("Error in updating owned building data: %w", err)
	}

	playerStats, err := GetPlayerInfoByID(ctx, playerID)
	if err != nil {
		return err
	}
	newSkill := playerStats.Skill_points + skillGain
	err = AddSkill(ctx, tx, playerID, newSkill)
	if err != nil {
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("Failed to commit to database: %w", err)
	}

	return nil
}

func CheckUpgrading(ctx context.Context, ownedBuildingID uuid.UUID) error {

	query := `SELECT upgrade_complete_at FROM owned_building WHERE id = $1`
	var finishTime *time.Time

	err := database.DB.QueryRow(ctx, query, ownedBuildingID).Scan(&finishTime)
	if err != nil {
		return fmt.Errorf("Error checking undergoing upgrades")
	}
	if finishTime != nil {
		if time.Now().Before(*finishTime) {
			return fmt.Errorf("Upgrade is still going on")
		}
	}

	return nil
}

func CollectResource(ctx context.Context, ownedBuildingID uuid.UUID) error {
	err := CheckUpgrading(ctx, ownedBuildingID)
	if err != nil {
		return err
	}

	var bDataID int
	var playerID uuid.UUID
	var lastCollect *time.Time

	query1 := `SELECT building_data_id, player_id, last_collected_at FROM owned_building WHERE id = $1`

	err = database.DB.QueryRow(ctx, query1, ownedBuildingID).Scan(&bDataID, &playerID, &lastCollect)
	if err != nil {
		return fmt.Errorf("Error in fetching last collected timestamp: %w", err)
	}

	bData, err := GetBuildingDataByID(ctx, bDataID)
	if err != nil {
		return err
	}

	resData, err := GetResBuildingData(ctx, bData.BuildingType, bData.BuildingLevel)
	if err != nil {
		return err
	}

	playerStats, err := GetPlayerInfoByID(ctx, playerID)
	if err != nil {
		return err
	}

	if lastCollect == nil {
		return fmt.Errorf("No resources generated yet")
	}

	tx, err := database.DB.Begin(ctx)
	if err != nil {
		return fmt.Errorf("Failed to begin databse transaction: %w", err)
	}

	defer tx.Rollback(ctx)

	timePassed := time.Since(*lastCollect)
	elixirGen := int(timePassed.Minutes()) * resData.ElixirGenPerMin
	pancakesGen := int(timePassed.Minutes()) * resData.PancakesGenPerMin

	newElixir := elixirGen + playerStats.Elixir
	newPancakes := pancakesGen + playerStats.Pancakes
	if elixirGen == 0 && pancakesGen == 0 {
		return fmt.Errorf("Nothing to cellect yet: no resources generated")
	}

	err = UpdateResources(ctx, tx, playerID, newPancakes, newElixir)
	if err != nil {
		return err
	}

	query2 := `UPDATE owned_building SET last_collected_at = $1 WHERE id = $2`
	_, err = tx.Exec(ctx, query2, time.Now(), ownedBuildingID)
	if err != nil {
		return fmt.Errorf("Error in updating timestamp in owned buildings: %w", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("Failed to commit to database: %w", err)
	}

	return nil
}

func GetPlayerStorageCapacity(ctx context.Context, playerID uuid.UUID) (int, int, error) {

	query := `
	SELECT bd.building_type, sd.max_storage FROM owned_building ob
	INNER JOIN building_data bd ON ob.building_data_id = bd.id
	INNER JOIN storage_building_data sd ON sd.building_data_id = bd.id
	WHERE ob.player_id = $1 AND ob.upgrade_complete_at = NULL
	`

	rows, err := database.DB.Query(ctx, query, playerID)
	if err != nil {
		return 0, 0, fmt.Errorf("Failed to fetch storage buildings: %w", err)
	}
	defer rows.Close()

	maxElixir := 1000
	maxPancakes := 1000

	for rows.Next() {

		var bType string
		var storage int

		err = rows.Scan(&bType, &storage)
		if err != nil {
			return 0, 0, err
		}
		if bType == "Elixir Storage" {
			maxElixir += storage
		}
		if bType == "Pancake Stack" {
			maxPancakes += storage
		}
	}

	return maxElixir, maxPancakes, nil
}
