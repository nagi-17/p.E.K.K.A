package models

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/nagi-17/p.E.K.K.A/internal/database"
)

func GetOwnedBuildingData(ctx context.Context, playerID uuid.UUID) ([]OwnedBuildingWithData, error) {
	checkQuery := `SELECT id FROM owned_building WHERE player_id = $1 AND upgrade_complete_at IS NOT NULL AND upgrade_complete_at <= NOW()`
	rowsCheck, err := database.DB.Query(ctx, checkQuery, playerID)
	if err == nil {
		var completedIDs []uuid.UUID
		for rowsCheck.Next() {
			var id uuid.UUID
			if errScan := rowsCheck.Scan(&id); errScan == nil {
				completedIDs = append(completedIDs, id)
			}
		}
		rowsCheck.Close()

		for _, bID := range completedIDs {
			_ = FinishUpgrade(ctx, bID)
		}
	}

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

func InsertBuilding(ctx context.Context, tx pgx.Tx, playerID uuid.UUID, bDataID int, x int, y int, upgrade_complete_at *time.Time, last_collected_at *time.Time) error {

	query := `
	INSERT INTO owned_building (id, player_id, building_data_id, pos_x, pos_y, upgrade_complete_at, last_collected_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	var id uuid.UUID = uuid.New()

	_, err := tx.Exec(ctx, query, id, playerID, bDataID, x, y, upgrade_complete_at, last_collected_at)
	if err != nil {
		return fmt.Errorf("Error in placing new building: %w", err)
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
		return fmt.Errorf("Failed to begin database transaction: %w", err)
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

func GetBuildingAndPlayerIDUsingOwnedBuildingID(ctx context.Context, ownedBuildingID uuid.UUID) (int, uuid.UUID, *time.Time, error) {
	query := `SELECT building_data_id, player_id, upgrade_complete_at FROM owned_building WHERE owned_building.id = $1 `

	var bDataID int
	var playerID uuid.UUID
	var curr_upgrade_complete_time *time.Time

	err := database.DB.QueryRow(ctx, query, ownedBuildingID).Scan(&bDataID, &playerID, &curr_upgrade_complete_time)
	if err != nil {
		return 0, uuid.Nil, nil, fmt.Errorf("Cannot fetch building id: %w", err)
	}

	return bDataID, playerID, curr_upgrade_complete_time, nil
}

func StartBuildingUpgradeQueryFunc(ctx context.Context, tx pgx.Tx, ownedBuildingID uuid.UUID, finishTime time.Time) error {

	query := `UPDATE owned_building SET upgrade_complete_at = $1 WHERE owned_building.id = $2`

	_, err := tx.Exec(ctx, query, finishTime, ownedBuildingID)
	if err != nil {
		return fmt.Errorf("Error in updating finish time of building: %w", err)
	}

	return nil
}

func FinishUpgrade(ctx context.Context, ownedBuildingID uuid.UUID) error {
	query1 := `
	SELECT ob.player_id, ob.upgrade_complete_at, bd.building_type, bd.building_level, bd.skill_on_upgrade, ob.last_collected_at
	FROM owned_building ob
	INNER JOIN building_data bd ON ob.building_data_id = bd.id
	WHERE ob.id = $1
	`
	var playerID uuid.UUID
	var upgradeCompleteAt *time.Time
	var bType string
	var currLevel int
	var skillGain int
	var lastCollectedAt *time.Time
	err := database.DB.QueryRow(ctx, query1, ownedBuildingID).Scan(&playerID, &upgradeCompleteAt, &bType, &currLevel, &skillGain, &lastCollectedAt)
	if err != nil {
		return fmt.Errorf("Failed to fetch building: %w", err)
	}

	if upgradeCompleteAt == nil {
		return fmt.Errorf("Building is not under upgrade")
	}
	if time.Now().Before(*upgradeCompleteAt) {
		return fmt.Errorf("Upgrade is still going on")
	}

	isInitialPlacement := lastCollectedAt != nil && lastCollectedAt.Unix() == 0

	upgradedLevel := currLevel
	if !isInitialPlacement {
		upgradedLevel = currLevel + 1
	}

	upgradedBData, err := GetBuildingDataByTypeLevel(ctx, bType, upgradedLevel)
	if err != nil {
		return fmt.Errorf("Failed to find data of upgraded building: %w", err)
	}

	tx, err := database.DB.Begin(ctx)
	if err != nil {
		return fmt.Errorf("Failed to begin database transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	var newLastCollected *time.Time = nil
	if isInitialPlacement {
		if bType == "Elixir Collector" || bType == "Pancake Machine" {
			nowTime := time.Now()
			newLastCollected = &nowTime
		}
	} else {
		newLastCollected = lastCollectedAt
	}

	query2 := `UPDATE owned_building SET building_data_id = $1, upgrade_complete_at = NULL, last_collected_at = $2 WHERE id = $3`
	_, err = tx.Exec(ctx, query2, upgradedBData.ID, newLastCollected, ownedBuildingID)
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
		return fmt.Errorf("Failed to begin database transaction: %w", err)
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
	WHERE ob.player_id = $1`

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

func CancelUpgrade(ctx context.Context, ownedBuildingID uuid.UUID) error {
	query1 := `
	SELECT ob.player_id, ob.upgrade_complete_at, bd.building_type, bd.building_level, bd.upgrade_cost_elixir, bd.upgrade_cost_pancakes, ob.last_collected_at
	FROM owned_building ob
	INNER JOIN building_data bd ON ob.building_data_id = bd.id
	WHERE ob.id = $1
	`
	var playerID uuid.UUID
	var upgradeCompleteAt *time.Time
	var bType string
	var currLevel int
	var costElixir int
	var costPancakes int
	var lastCollectedAt *time.Time

	err := database.DB.QueryRow(ctx, query1, ownedBuildingID).Scan(&playerID, &upgradeCompleteAt, &bType, &currLevel, &costElixir, &costPancakes, &lastCollectedAt)
	if err != nil {
		return fmt.Errorf("Failed to fetch building: %w", err)
	}

	if upgradeCompleteAt == nil {
		return fmt.Errorf("Building is not under upgrade")
	}

	if time.Now().After(*upgradeCompleteAt) {
		return fmt.Errorf("Upgrade has already finished")
	}

	tx, err := database.DB.Begin(ctx)
	if err != nil {
		return fmt.Errorf("Failed to begin database transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	isInitialPlacement := lastCollectedAt != nil && lastCollectedAt.Unix() == 0

	if isInitialPlacement {
		queryDelete := `DELETE FROM owned_building WHERE id = $1`
		_, err = tx.Exec(ctx, queryDelete, ownedBuildingID)
		if err != nil {
			return fmt.Errorf("Error in deleting cancelled building placement: %w", err)
		}
	} else {
		queryRevert := `UPDATE owned_building SET upgrade_complete_at = NULL WHERE id = $1`
		_, err = tx.Exec(ctx, queryRevert, ownedBuildingID)
		if err != nil {
			return fmt.Errorf("Error in reverting building upgrade: %w", err)
		}
	}

	refundElixir := costElixir / 2
	refundPancakes := costPancakes / 2

	playerStats, err := GetPlayerInfoByID(ctx, playerID)
	if err != nil {
		return err
	}

	newPancakes := playerStats.Pancakes + refundPancakes
	newElixir := playerStats.Elixir + refundElixir

	err = UpdateResources(ctx, tx, playerID, newPancakes, newElixir)
	if err != nil {
		return fmt.Errorf("Error updating player resources: %w", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("Failed to commit database transaction: %w", err)
	}

	return nil
}
