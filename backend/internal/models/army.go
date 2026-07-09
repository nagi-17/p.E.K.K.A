package models

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/nagi-17/p.E.K.K.A/internal/database"
	"github.com/nagi-17/p.E.K.K.A/internal/utils"
)

func TrainArmy(ctx context.Context, playerID string, troopType string, quantity int) error {
	playerUUID, err := uuid.Parse(playerID)
	if err != nil {
		return fmt.Errorf("Error in parsing player id: %w", err)
	}

	playerTroopInfo, err := GetPlayerTroopLevel(ctx, playerID, troopType)
	if err != nil {
		return fmt.Errorf("Error fetching player troop level: %w", err)
	}

	var trainTroopID, space int
	query1 := `SELECT id, space_occupied_in_army FROM troop_data WHERE troop_type = $1 AND troop_level = $2`

	err = database.DB.QueryRow(ctx, query1, troopType, playerTroopInfo.CurrentLevel).Scan(&trainTroopID, &space)
	if err != nil {
		return fmt.Errorf("Error fetching specified troop stats: %w", err)
	}

	query2 := `SELECT troop_data_id, quantity FROM trained_troop WHERE player_id = $1`
	rows, err := database.DB.Query(ctx, query2, playerUUID)
	if err != nil {
		return fmt.Errorf("Error fetching trained army data: %w", err)
	}
	defer rows.Close()

	var currSpace int
	var allTrainedTroops []TrainedTroop
	for rows.Next() {
		var t TrainedTroop

		err = rows.Scan(&t.TroopDataID, &t.Quantity)
		if err != nil {
			return fmt.Errorf("Error in scanning trained troop data/rows: %w", err)
		}

		allTrainedTroops = append(allTrainedTroops, t)
		troopData, err := GetTroopStats(ctx, t.TroopDataID)
		if err == nil {
			currSpace += t.Quantity * troopData.SpaceOccupiedInArmy
		}
	}

	var maxSpace int
	allOwnedBuildings, err := GetOwnedBuildingData(ctx, playerUUID)

	if err != nil {
		return fmt.Errorf("Error in fetching army camps housing data: %w", err)
	}

	for i := 0; i < len(allOwnedBuildings); i++ {
		if allOwnedBuildings[i].BuildingType == "Army Camp" {

			if CheckUpgrading(ctx, allOwnedBuildings[i].ID) != nil {
				continue
			}

			query3 := `SELECT housing_space FROM army_camp_data WHERE building_data_id = $1`
			var tempSpace int

			err = database.DB.QueryRow(ctx, query3, allOwnedBuildings[i].BuildingDataID).Scan(&tempSpace)
			if err == nil {
				maxSpace += tempSpace
			}
		}
	}

	if maxSpace == 0 {
		return fmt.Errorf("All army camps are under upgrade")
	}

	spaceReq := quantity * space
	if currSpace+spaceReq > maxSpace {
		return utils.NewUserError("Insufficient housing space")
	}

	for i := 0; i < len(allTrainedTroops); i++ {
		if allTrainedTroops[i].TroopDataID == trainTroopID {
			newQuantity := quantity + allTrainedTroops[i].Quantity
			query4 := `UPDATE trained_troop SET quantity = $1 WHERE player_id = $2 AND troop_data_id = $3`

			_, err = database.DB.Exec(ctx, query4, newQuantity, playerUUID, trainTroopID)
			if err != nil {
				return fmt.Errorf("Failed to update troops in army: %w", err)
			}
			return nil
		}
	}

	var id uuid.UUID = uuid.New()
	query5 := `INSERT INTO trained_troop (id, player_id, troop_data_id, quantity) VALUES ($1, $2, $3, $4)`

	_, err = database.DB.Exec(ctx, query5, id, playerUUID, trainTroopID, quantity)
	if err != nil {
		return fmt.Errorf("Failed to add troops to army: %w", err)
	}

	return nil
}

type TrainedTroopStatus struct {
	TroopType         string     `json:"troop_type"`
	CurrentLevel      int        `json:"level"`
	Quantity          int        `json:"quantity"`
	UpgradeCompleteAt *time.Time `json:"upgrade_complete_at"`
}

func GetTrainedArmy(ctx context.Context, playerID string) ([]TrainedTroopStatus, error) {
	playerUUID, err := uuid.Parse(playerID)
	if err != nil {
		return nil, fmt.Errorf("Error in parsing player id: %w", err)
	}

	checkQuery := `SELECT troop_type FROM player_troop_level WHERE player_id = $1 AND upgrade_complete_at IS NOT NULL AND upgrade_complete_at <= NOW()`
	rowsCheck, err := database.DB.Query(ctx, checkQuery, playerUUID)
	if err == nil {
		var completedTypes []string
		for rowsCheck.Next() {
			var tType string
			if errScan := rowsCheck.Scan(&tType); errScan == nil {
				completedTypes = append(completedTypes, tType)
			}
		}
		rowsCheck.Close()

		for _, tType := range completedTypes {
			_ = FinishTroopUpgrade(ctx, playerID, tType)
		}
	}

	query := `
		SELECT ptl.troop_type, ptl.current_level, ptl.upgrade_complete_at, COALESCE(tt.quantity, 0) as quantity
		FROM player_troop_level ptl
		LEFT JOIN troop_data td ON td.troop_type = ptl.troop_type AND td.troop_level = ptl.current_level
		LEFT JOIN trained_troop tt ON tt.player_id = ptl.player_id AND tt.troop_data_id = td.id
		WHERE ptl.player_id = $1
		ORDER BY ptl.troop_type;
	`
	rows, err := database.DB.Query(ctx, query, playerUUID)
	if err != nil {
		return nil, fmt.Errorf("Error querying trained army: %w", err)
	}
	defer rows.Close()

	var army []TrainedTroopStatus
	for rows.Next() {
		var t TrainedTroopStatus
		err = rows.Scan(&t.TroopType, &t.CurrentLevel, &t.UpgradeCompleteAt, &t.Quantity)
		if err != nil {
			return nil, fmt.Errorf("Error scanning trained troop status: %w", err)
		}
		army = append(army, t)
	}
	return army, nil
}

func DiscardTroops(ctx context.Context, playerID string, troopType string, quantity int) error {
	playerUUID, err := uuid.Parse(playerID)
	if err != nil {
		return fmt.Errorf("Error in parsing player id: %w", err)
	}

	playerTroopInfo, err := GetPlayerTroopLevel(ctx, playerID, troopType)
	if err != nil {
		return fmt.Errorf("Error fetching player troop level: %w", err)
	}

	var troopDataID int
	query1 := `SELECT id FROM troop_data WHERE troop_type = $1 AND troop_level = $2`
	err = database.DB.QueryRow(ctx, query1, troopType, playerTroopInfo.CurrentLevel).Scan(&troopDataID)
	if err != nil {
		return fmt.Errorf("Error fetching specified troop stats: %w", err)
	}

	var currentQuantity int
	query2 := `SELECT quantity FROM trained_troop WHERE player_id = $1 AND troop_data_id = $2`
	err = database.DB.QueryRow(ctx, query2, playerUUID, troopDataID).Scan(&currentQuantity)
	if err != nil {
		return utils.NewUserError("No troops of this type trained")
	}

	if currentQuantity < quantity {
		return utils.NewUserError("Cannot discard more troops than you have trained")
	}

	newQuantity := currentQuantity - quantity
	if newQuantity == 0 {
		query3 := `DELETE FROM trained_troop WHERE player_id = $1 AND troop_data_id = $2`
		_, err = database.DB.Exec(ctx, query3, playerUUID, troopDataID)
	} else {
		query3 := `UPDATE trained_troop SET quantity = $1 WHERE player_id = $2 AND troop_data_id = $3`
		_, err = database.DB.Exec(ctx, query3, newQuantity, playerUUID, troopDataID)
	}

	if err != nil {
		return fmt.Errorf("Failed to discard troops: %w", err)
	}

	return nil
}
