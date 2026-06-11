package models

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/nagi-17/p.E.K.K.A/internal/database"
)

type BuildingData struct {
	ID                   int           `db:"id"`
	BuildingType         string        `db:"building_type"`
	BuildingLevel        int           `db:"building_level"`
	Health               int           `db:"health"`
	Width                int           `db:"width"`
	Height               int           `db:"height"`
	BuildTime            time.Duration `db:"build_time"`
	UpgradeCostElixir    int           `db:"upgrade_cost_elixir"`
	UpgradeCostPancakes  int           `db:"upgrade_cost_pancakes"`
	UpgradeTime          time.Duration `db:"upgrade_time"`
	MaxQuantityAvailable int           `db:"max_quantity_available"`
	SkillOnUpgrade       int           `db:"skill_on_upgrade"`
}

type OwnedBuildingData struct {
	ID                uuid.UUID  `db:"id"`
	PlayerID          uuid.UUID  `db:"player_id"`
	BuildingDataID    int        `db:"building_data_id"`
	PosX              int        `db:"pos_x"`
	PosY              int        `db:"pos_y"`
	UpgradeCompleteAt *time.Time `db:"upgrade_complete_at"`
	LastCollectedAt   *time.Time `db:"last_collected_at"`
}

type TownHallData struct {
	BuildingData
	MinSkillPointsBeforeUpgrade int `db:"min_skill_points_before_upgrade"`
}

type DefenseBuildingData struct {
	BuildingData
	BuildingRange           int `db:"building_range"`
	DamagePerSec            int `db:"damage_per_sec"`
	DamagePerShot           int `db:"damage_per_shot"`
	MaxPossibleUpgradeLevel int `db:"max_possible_upgrade_level"`
	UnlockTownHallLevel     int `db:"unlock_town_hall_level"`
}

type ResourceBuildingData struct {
	BuildingData
	ElixirGenPerMin         int `db:"elixir_gen_per_min"`
	PancakesGenPerMin       int `db:"pancakes_gen_per_min"`
	MaxPossibleUpgradeLevel int `db:"max_possible_upgrade_level"`
	UnlockTownHallLevel     int `db:"unlock_town_hall_level"`
}

type StorageBuildingData struct {
	BuildingData
	MaxStorage              int `db:"max_storage"`
	MaxPossibleUpgradeLevel int `db:"max_possible_upgrade_level"`
	UnlockTownHallLevel     int `db:"unlock_town_hall_level"`
}

type LaboratoryData struct {
	BuildingData
	MaxTroopUpgradeLevel    int `db:"max_troop_upgrade_level"`
	MaxPossibleUpgradeLevel int `db:"max_possible_upgrade_level"`
	UnlockTownHallLevel     int `db:"unlock_town_hall_level"`
}

type OwnedBuildingWithData struct {
	OwnedBuildingData
	BuildingType         string        `db:"building_type"`
	BuildingLevel        int           `db:"building_level"`
	Health               int           `db:"health"`
	Width                int           `db:"width"`
	Height               int           `db:"height"`
	BuildTime            time.Duration `db:"build_time"`
	UpgradeCostElixir    int           `db:"upgrade_cost_elixir"`
	UpgradeCostPancakes  int           `db:"upgrade_cost_pancakes"`
	UpgradeTime          time.Duration `db:"upgrade_time"`
	MaxQuantityAvailable int           `db:"max_quantity_available"`
	SkillOnUpgrade       int           `db:"skill_on_upgrade"`
}

func GetBuildingDataByID(ctx context.Context, buildingID int) (*BuildingData, error) {
	query := `
	SELECT id, building_type, building_level, health, width, height,
	build_time, upgrade_cost_elixir, upgrade_cost_pancakes, upgrade_time,
	max_quantity_available, skill_on_upgrade
	FROM building_data
	WHERE id = $1
	`

	var bData BuildingData
	err := database.DB.QueryRow(ctx, query, buildingID).Scan(&bData.ID, &bData.BuildingType, &bData.BuildingLevel, &bData.Health,
		&bData.Width, &bData.Height, &bData.BuildTime, &bData.UpgradeCostElixir,
		&bData.UpgradeCostPancakes, &bData.UpgradeTime, &bData.MaxQuantityAvailable,
		&bData.SkillOnUpgrade)
	if err != nil {
		return nil, fmt.Errorf("Error in fetching building (static)data: %w", err)
	}
	return &bData, nil
}

func GetTownHallData(ctx context.Context, level int) (*TownHallData, error) {
	query := `
	SELECT b.id, b.building_type, b.building_level, b.health,
	b.width, b.height, b.build_time, b.upgrade_cost_elixir,
	b.upgrade_cost_pancakes, b.upgrade_time, b.max_quantity_available,
	b.skill_on_upgrade, t.min_skill_points_before_upgrade
	FROM building_data b
	JOIN town_hall_data t ON t.building_data_id = b.id
	WHERE b.building_level = $1
	`
	var thData TownHallData
	err := database.DB.QueryRow(ctx, query, level).Scan(&thData.BuildingData.ID, &thData.BuildingData.BuildingType, &thData.BuildingData.BuildingLevel,
		&thData.BuildingData.Health, &thData.BuildingData.Width, &thData.BuildingData.Height, &thData.BuildingData.BuildTime, &thData.BuildingData.UpgradeCostElixir,
		&thData.BuildingData.UpgradeCostPancakes, &thData.BuildingData.UpgradeTime, &thData.BuildingData.MaxQuantityAvailable,
		&thData.BuildingData.SkillOnUpgrade, &thData.MinSkillPointsBeforeUpgrade)
	if err != nil {
		return nil, fmt.Errorf("Error in fetching townhall data: %w", err)
	}
	return &thData, nil
}

func GetDefBuildingData(ctx context.Context, bType string, level int) (*DefenseBuildingData, error) {
	query := `
	SELECT b.id, b.building_type, b.building_level, b.health,
	b.width, b.height, b.build_time, b.upgrade_cost_elixir,
	b.upgrade_cost_pancakes, b.upgrade_time, b.max_quantity_available,
	b.skill_on_upgrade, d.building_range, d.damage_per_sec, d.damage_per_shot,
	d.max_possible_upgrade_level, d.unlock_town_hall_level
	FROM building_data b
	JOIN defense_building_data d ON d.building_data_id = b.id
	WHERE b.building_type = $1 AND b.building_level = $2
	`
	var defData DefenseBuildingData
	err := database.DB.QueryRow(ctx, query, bType, level).Scan(&defData.BuildingData.ID, &defData.BuildingData.BuildingType, &defData.BuildingData.BuildingLevel,
		&defData.BuildingData.Health, &defData.BuildingData.Width, &defData.BuildingData.Height, &defData.BuildingData.BuildTime, &defData.BuildingData.UpgradeCostElixir,
		&defData.BuildingData.UpgradeCostPancakes, &defData.BuildingData.UpgradeTime, &defData.BuildingData.MaxQuantityAvailable, &defData.BuildingData.SkillOnUpgrade,
		&defData.BuildingRange, &defData.DamagePerSec, &defData.DamagePerShot, &defData.MaxPossibleUpgradeLevel, &defData.UnlockTownHallLevel)
	if err != nil {
		return nil, fmt.Errorf("Error in fetching def. data: %w", err)
	}
	return &defData, nil
}

func GetResBuildingData(ctx context.Context, bType string, level int) (*ResourceBuildingData, error) {
	query := `
	SELECT b.id, b.building_type, b.building_level, b.health,
	b.width, b.height, b.build_time, b.upgrade_cost_elixir,
	b.upgrade_cost_pancakes, b.upgrade_time, b.max_quantity_available,
	b.skill_on_upgrade, r.elixir_gen_per_min, r.pancakes_gen_per_min,
	r.max_possible_upgrade_level, r.unlock_town_hall_level
	FROM building_data b
	JOIN resource_building_data r ON r.building_data_id = b.id
	WHERE b.building_type = $1 AND b.building_level = $2
	`
	var resData ResourceBuildingData
	err := database.DB.QueryRow(ctx, query, bType, level).Scan(&resData.BuildingData.ID, &resData.BuildingData.BuildingType, &resData.BuildingData.BuildingLevel,
		&resData.BuildingData.Health, &resData.BuildingData.Width, &resData.BuildingData.Height, &resData.BuildingData.BuildTime, &resData.BuildingData.UpgradeCostElixir,
		&resData.BuildingData.UpgradeCostPancakes, &resData.BuildingData.UpgradeTime, &resData.BuildingData.MaxQuantityAvailable, &resData.BuildingData.SkillOnUpgrade,
		&resData.ElixirGenPerMin, &resData.PancakesGenPerMin, &resData.MaxPossibleUpgradeLevel, &resData.UnlockTownHallLevel)
	if err != nil {
		return nil, fmt.Errorf("Error in fetching res. data: %w", err)
	}
	return &resData, nil
}

func GetStrgBuildingData(ctx context.Context, bType string, level int) (*StorageBuildingData, error) {
	query := `
	SELECT b.id, b.building_type, b.building_level, b.health,
	b.width, b.height, b.build_time, b.upgrade_cost_elixir,
	b.upgrade_cost_pancakes, b.upgrade_time, b.max_quantity_available,
	b.skill_on_upgrade, s.max_storage,
	s.max_possible_upgrade_level, s.unlock_town_hall_level
	FROM building_data b
	JOIN storage_building_data s ON s.building_data_id = b.id
	WHERE b.building_type = $1 AND b.building_level = $2
	`
	var strgData StorageBuildingData
	err := database.DB.QueryRow(ctx, query, bType, level).Scan(&strgData.BuildingData.ID, &strgData.BuildingData.BuildingType, &strgData.BuildingData.BuildingLevel,
		&strgData.BuildingData.Health, &strgData.BuildingData.Width, &strgData.BuildingData.Height, &strgData.BuildingData.BuildTime, &strgData.BuildingData.UpgradeCostElixir,
		&strgData.BuildingData.UpgradeCostPancakes, &strgData.BuildingData.UpgradeTime, &strgData.BuildingData.MaxQuantityAvailable, &strgData.BuildingData.SkillOnUpgrade,
		&strgData.MaxStorage, &strgData.MaxPossibleUpgradeLevel, &strgData.UnlockTownHallLevel)
	if err != nil {
		return nil, fmt.Errorf("Error in fetching storage data: %w", err)
	}
	return &strgData, nil
}

func GetLabData(ctx context.Context, bType string, level int) (*LaboratoryData, error) {
	query := `
	SELECT b.id, b.building_type, b.building_level, b.health,
	b.width, b.height, b.build_time, b.upgrade_cost_elixir,
	b.upgrade_cost_pancakes, b.upgrade_time, b.max_quantity_available,
	b.skill_on_upgrade, l.max_troop_upgrade_level,
	l.max_possible_upgrade_level, l.unlock_town_hall_level
	FROM building_data b
	JOIN laboratory_data l ON l.building_data_id = b.id
	WHERE b.building_type = $1 AND b.building_level = $2
	`
	var labData LaboratoryData
	err := database.DB.QueryRow(ctx, query, bType, level).Scan(&labData.BuildingData.ID, &labData.BuildingData.BuildingType, &labData.BuildingData.BuildingLevel,
		&labData.BuildingData.Health, &labData.BuildingData.Width, &labData.BuildingData.Height, &labData.BuildingData.BuildTime, &labData.BuildingData.UpgradeCostElixir,
		&labData.BuildingData.UpgradeCostPancakes, &labData.BuildingData.UpgradeTime, &labData.BuildingData.MaxQuantityAvailable, &labData.BuildingData.SkillOnUpgrade,
		&labData.MaxTroopUpgradeLevel, &labData.MaxPossibleUpgradeLevel, &labData.UnlockTownHallLevel)
	if err != nil {
		return nil, fmt.Errorf("Error in fetching lab data: %w", err)
	}
	return &labData, nil
}
