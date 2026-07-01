package models

import (
	"time"

	"github.com/google/uuid"
)

type BattleLog struct {
	ID             uuid.UUID `db:"id" json:"id"`
	AttackerID     uuid.UUID `db:"attacker_id" json:"attacker_id"`
	DefenderID     uuid.UUID `db:"defender_id" json:"defender_id"`
	ElixirLooted   int       `db:"elixir_looted" json:"elixir_looted"`
	PancakesLooted int       `db:"pancakes_looted" json:"pancakes_looted"`
	DamagePercent  float32   `db:"damage_percent" json:"damage_percent"`
	TimeOfBattle   time.Time `db:"time_of_battle" json:"time_of_battle"`
}

type AttackingTroop struct {
	TroopType  string
	Health     int
	DPS        int
	DPA        int
	TroopRange int
	Quantity   int
	Airborne   bool
}

type DefenseSnapshot struct {
	BuildingType  string
	BuildingLevel int
	Health        int
	BuildingRange int
	DamagePerSec  int
	DamagePerShot int
	PosX          int
	PosY          int
	Width         int
	Height        int
}

type OwnedBuildingWithData struct {
	OwnedBuildingData
	BuildingType         string        `db:"building_type" json:"building_type"`
	BuildingLevel        int           `db:"building_level" json:"building_level"`
	Health               int           `db:"health" json:"health"`
	Width                int           `db:"width" json:"width"`
	Height               int           `db:"height" json:"height"`
	BuildTime            time.Duration `db:"build_time" json:"build_time"`
	UpgradeCostElixir    int           `db:"upgrade_cost_elixir" json:"upgrade_cost_elixir"`
	UpgradeCostPancakes  int           `db:"upgrade_cost_pancakes" json:"upgrade_cost_pancakes"`
	UpgradeTime          time.Duration `db:"upgrade_time" json:"upgrade_time"`
	MaxQuantityAvailable int           `db:"max_quantity_available" json:"max_quantity_available"`
	SkillOnUpgrade       int           `db:"skill_on_upgrade" json:"skill_on_upgrade"`
}

type OwnedBuildingData struct {
	ID                uuid.UUID  `db:"id" json:"id"`
	PlayerID          uuid.UUID  `db:"player_id" json:"player_id"`
	BuildingDataID    int        `db:"building_data_id" json:"building_data_id"`
	PosX              int        `db:"pos_x" json:"pos_x"`
	PosY              int        `db:"pos_y" json:"pos_y"`
	UpgradeCompleteAt *time.Time `db:"upgrade_complete_at" json:"upgrade_complete_at"`
	LastCollectedAt   *time.Time `db:"last_collected_at" json:"last_collected_at"`
}

type TrainedTroop struct {
	ID          uuid.UUID `db:"id"`
	PlayerID    string    `db:"player_id"`
	TroopDataID int       `db:"troop_data_id"`
	Quantity    int       `db:"quantity"`
}
