package models

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/nagi-17/p.E.K.K.A/internal/database"
)

type BattleLog struct {
	ID             uuid.UUID `db:"id"`
	AttackerID     uuid.UUID `db:"attacker_id"`
	DefenderID     uuid.UUID `db:"defender_id"`
	ElixirLooted   int       `db:"elixir_looted"`
	PancakesLooted int       `db:"pancakes_looted"`
	DamagePercent  float32   `db:"damage_percent"`
	TimeOfBattle   time.Time `db:"time_of_battle"`
}

type OpponentData struct {
	PlayerID    uuid.UUID `json:"player_id"`
	Trophies    int       `json:"trophies"`
	SkillPoints int       `json:"skill_points"`
	Elixir      int       `json:"elixir"`
	Pancakes    int       `json:"pancakes"`
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
}

func FindOpponent(ctx context.Context, attackerID string) (*OpponentData, error) {
	attackerUUID, err := uuid.Parse(attackerID)
	if err != nil {
		return nil, fmt.Errorf("Can't parse attacker id")
	}

	attackerStats, err := GetPlayerInfoByID(ctx, attackerUUID)
	if err != nil {
		return nil, fmt.Errorf("Failed to load attacker stats: %w", err)
	}

	maxSkill := attackerStats.Skill_points + 100
	minSkill := attackerStats.Skill_points - 100
	if minSkill < 0 {
		minSkill = 0
	}

	query := `SELECT player_id, trophies, skill_points, elixir, pancakes FROM player_info 
	WHERE player_id != $1 AND (shield_end_time IS NULL OR shield_end_time < NOW()) AND skill_points BETWEEN $2 AND $3 ORDER BY RANDOM() LIMIT 1`

	var defender OpponentData
	err = database.DB.QueryRow(ctx, query, attackerUUID, minSkill, maxSkill).Scan(&defender.PlayerID, &defender.Trophies, &defender.SkillPoints, &defender.Elixir, &defender.Pancakes)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, fmt.Errorf("Haha...SKILL ISSUE: No base found in your skill range")
		}
		return nil, fmt.Errorf("Failed to fetch opponent base: %w", err)
	}

	return &defender, nil
}

func GetAttackingArmy(ctx context.Context, playerID uuid.UUID) ([]AttackingTroop, error) {
	query := `SELECT td.troop_type, td.health, td.dps, td.dpa, td.troop_range, td.airborne, tt.quantity
	FROM trained_troop tt 
	INNER JOIN troop_data td ON tt.troop_data_id = td.id 
	WHERE tt.player_id = $1 AND tt.quantity > 0`

	rows, err := database.DB.Query(ctx, query, playerID)
	if err != nil {
		return nil, fmt.Errorf("Error in fetching attacking army: %w", err)
	}
	defer rows.Close()

	var army []AttackingTroop
	for rows.Next() {
		var t AttackingTroop
		err = rows.Scan(&t.TroopType, &t.Health, &t.DPS, &t.DPA, &t.TroopRange, &t.Airborne, &t.Quantity)
		if err != nil {
			return nil, fmt.Errorf("Error scanning army rows: %w", err)
		}
		army = append(army, t)
	}

	if len(army) == 0 {
		return nil, fmt.Errorf("You have no troops to attack with")
	}
	return army, nil
}

func GetDefenseSnapshot(ctx context.Context, playerID uuid.UUID) ([]DefenseSnapshot, error) {
	query := `
	SELECT bd.building_type, bd.building_level, bd.health, dd.building_range, dd.damage_per_sec, dd.damage_per_shot, ob.pos_x, ob.pos_y
	FROM owned_building ob 
	INNER JOIN building_data bd ON ob.building_data_id = bd.id 
	INNER JOIN defense_building_data dd ON dd.building_data_id = bd.id 
	WHERE ob.player_id = $1`

	rows, err := database.DB.Query(ctx, query, playerID)
	if err != nil {
		return nil, fmt.Errorf("Error in fetching defense snapshot: %w", err)
	}
	defer rows.Close()

	var defenses []DefenseSnapshot
	for rows.Next() {
		var d DefenseSnapshot
		err = rows.Scan(&d.BuildingType, &d.BuildingLevel, &d.Health, &d.BuildingRange, &d.DamagePerSec, &d.DamagePerShot, &d.PosX, &d.PosY)
		if err != nil {
			return nil, fmt.Errorf("Error in scanning defense snapshot rows: %w", err)
		}
		defenses = append(defenses, d)
	}

	return defenses, nil
}

func CreateBattleLog(ctx context.Context, tx pgx.Tx, attackerID uuid.UUID, defenderID uuid.UUID, elixirLooted int, pancakesLooted int, damagePercent float32) error {
	query := `
	INSERT INTO battle_log (id, attacker_id, defender_id, elixir_looted, pancakes_looted, damage_percent, time_of_battle)
	VALUES ($1, $2, $3, $4, $5, $6, $7)`

	id := uuid.New()
	_, err := tx.Exec(ctx, query, id, attackerID, defenderID, elixirLooted, pancakesLooted, damagePercent, time.Now())
	if err != nil {
		return fmt.Errorf("Error in creating battle log: %w", err)
	}

	return nil
}
