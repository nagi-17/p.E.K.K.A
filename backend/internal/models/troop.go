package models

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/nagi-17/p.E.K.K.A/internal/database"
)

type TroopData struct {
	ID                  int    `db:"id"`
	TroopType           string `db:"troop_type"`
	TroopLevel          int    `db:"troop_level"`
	Health              int    `db:"health"`
	DPS                 int    `db:"dps"`
	DPA                 int    `db:"dpa"`
	TroopRange          int    `db:"troop_range"`
	UpgradeCostElixir   int    `db:"upgrade_cost_elixir"`
	SpaceOccupiedInArmy int    `db:"space_occupied_in_army"`
	MovSpeed            int    `db:"mov_speed"`
	AttackSpeed         int    `db:"attack_speed"`
	LabLevelReq         int    `db:"lab_level_req"`
	Airborne            bool   `db:"airborne"`
	UpgradeTime         int    `db:"upgrade_time"`
}

type TrainedTroop struct {
	ID          uuid.UUID `db:"id"`
	PlayerID    string    `db:"player_id"`
	TroopDataID int       `db:"troop_data_id"`
	Quantity    int       `db:"quantity"`
}

type PlayerTroopLevel struct {
	ID                uuid.UUID  `db:"id"`
	PlayerID          uuid.UUID  `db:"player_id"`
	TroopType         string     `db:"troop_type"`
	CurrentLevel      int        `db:"current_level"`
	UpgradeCompleteAt *time.Time `db:"upgrade_complete_at"`
}

func GetTroopStats(ctx context.Context, troopID int) (*TroopData, error) {

	query := `SELECT id, troop_type, troop_level, health, dps, dpa, troop_range,
	upgrade_cost_elixir, space_occupied_in_army, mov_speed, attack_speed,
	lab_level_req, airborne, upgrade_time FROM troop_data WHERE id = $1`

	var x TroopData
	err := database.DB.QueryRow(ctx, query, troopID).Scan(&x.ID, &x.TroopType,
		&x.TroopLevel, &x.Health, &x.DPS, &x.DPA, &x.TroopRange, &x.UpgradeCostElixir,
		&x.SpaceOccupiedInArmy, &x.MovSpeed, &x.AttackSpeed, &x.LabLevelReq, &x.Airborne, &x.UpgradeTime)
	if err != nil {
		return nil, fmt.Errorf("Error fetching troop (static)data")
	}

	return &x, nil
}

func GetPlayerTroopLevel(ctx context.Context, playerID string, troopType string) (*PlayerTroopLevel, error) {

	playerUUID, err := uuid.Parse(playerID)
	if err != nil {
		return nil, fmt.Errorf("Error in parsing player id")
	}

	var x PlayerTroopLevel

	query := `SELECT id, player_id, troop_type, current_level, upgrade_complete_at
	FROM player_troop_level WHERE player_id = $1 AND troop_type = $2`
	err = database.DB.QueryRow(ctx, query, playerUUID, troopType).Scan(&x.ID, &x.PlayerID,
		&x.TroopType, &x.CurrentLevel, &x.UpgradeCompleteAt)
	if err != nil {
		return nil, fmt.Errorf("Error in fetching troop level: %w", err)
	}

	return &x, nil
}

func StartTroopUpgrade(ctx context.Context, playerID string, troopType string) error {

	playerUUID, err := uuid.Parse(playerID)
	if err != nil {
		return fmt.Errorf("Error in parsing player id")
	}

	playerTroopInfo, err := GetPlayerTroopLevel(ctx, playerID, troopType)
	if err != nil {
		return fmt.Errorf("Error in fetching player troop data: %w", err)
	}

	if playerTroopInfo.UpgradeCompleteAt != nil && !playerTroopInfo.UpgradeCompleteAt.After(time.Now()) {
		err = FinishTroopUpgrade(ctx, playerID, troopType)
		if err != nil {
			return err
		}
	}

	if playerTroopInfo.UpgradeCompleteAt != nil && playerTroopInfo.UpgradeCompleteAt.After(time.Now()) {
		return fmt.Errorf("Upgrade is already in progress")
	}
	if playerTroopInfo.CurrentLevel == 4 {
		return fmt.Errorf("Troop is maxed out")
	}

	var labLevel int
	buildings, err := GetOwnedBuildingData(ctx, playerUUID)
	for i := 0; i < len(buildings); i++ {
		if buildings[i].BuildingType == "Laboratory" {
			err = CheckUpgrading(ctx, buildings[i].OwnedBuildingData.ID)
			if err != nil {
				return fmt.Errorf("Can't upgrade troop : Laboratory is under upgrade")
			}
			labLevel = buildings[i].BuildingLevel
			break
		}
	}
	if labLevel == 0 {
		return fmt.Errorf("Can't upgrade troops: build laboratory")
	}

	var labReq int
	nextLevel := playerTroopInfo.CurrentLevel + 1
	var cost int
	var upgradeTime int

	query1 := `SELECT lab_level_req, upgrade_cost_elixir, upgrade_time FROM troop_data WHERE troop_type = $1 AND troop_level = $2`
	err = database.DB.QueryRow(ctx, query1, troopType, nextLevel).Scan(&labReq, &cost, &upgradeTime)
	if err != nil {
		return fmt.Errorf("Can't fetch lab data: %w", err)
	}

	if labLevel < labReq {
		return fmt.Errorf("%d level Lab req. to upgrade this troop", labReq)
	}

	playerStats, err := GetPlayerInfoByID(ctx, playerUUID)
	if err != nil {
		return fmt.Errorf("Error in fetching player stats: %w", err)
	}

	if cost > playerStats.Elixir {
		return fmt.Errorf("Insufficient elixir")
	}

	tx, err := database.DB.Begin(ctx)
	if err != nil {
		return fmt.Errorf("Failed to begin databse transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	finishTime := time.Now().Add(time.Duration(upgradeTime) * time.Second)
	query2 := `UPDATE player_troop_level SET upgrade_complete_at = $1 WHERE player_id = $2 AND troop_type = $3`
	_, err = tx.Exec(ctx, query2, finishTime, playerUUID, troopType)
	if err != nil {
		return fmt.Errorf("Error in updating databse with update time: %w", err)
	}

	newElixir := playerStats.Elixir - cost
	err = UpdateResources(ctx, tx, playerUUID, playerStats.Pancakes, newElixir)
	if err != nil {
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("Failed to commit to database: %w", err)
	}

	return nil
}

func FinishTroopUpgrade(ctx context.Context, playerID string, troopType string) error {
	playerUUID, err := uuid.Parse(playerID)
	if err != nil {
		return fmt.Errorf("Error in parsing player id")
	}

	playerTroopInfo, err := GetPlayerTroopLevel(ctx, playerID, troopType)
	if err != nil {
		return fmt.Errorf("Error fetching player troop level data: %w", err)
	}

	if playerTroopInfo.UpgradeCompleteAt == nil {
		return fmt.Errorf("No upgrage is currently in progress")
	}

	if playerTroopInfo.UpgradeCompleteAt.After(time.Now()) {
		timeLeft := time.Until(*playerTroopInfo.UpgradeCompleteAt)
		return fmt.Errorf("Upgrade is still under process - time remaining: %v", timeLeft.Round(time.Second))
	}

	nextLevel := playerTroopInfo.CurrentLevel + 1
	query := `UPDATE player_troop_level SET current_level = $1, upgrade_complete_at = NULL WHERE player_id = $2 AND troop_type = $3`

	_, err = database.DB.Exec(ctx, query, nextLevel, playerUUID, troopType)
	if err != nil {
		return fmt.Errorf("Failed to finish troop upgrade: %w", err)
	}

	return nil
}
