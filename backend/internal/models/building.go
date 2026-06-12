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
	_, err = tx.Exec(ctx, query, id, playerID, bData.ID, x, y, upgrade_complete_at, last_collected_at)
	if err != nil {
		return fmt.Errorf("Error in placing new building: %w", err)
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
