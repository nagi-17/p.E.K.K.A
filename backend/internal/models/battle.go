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
	ID             uuid.UUID `db:"id" json:"id"`
	AttackerID     uuid.UUID `db:"attacker_id" json:"attacker_id"`
	DefenderID     uuid.UUID `db:"defender_id" json:"defender_id"`
	ElixirLooted   int       `db:"elixir_looted" json:"elixir_looted"`
	PancakesLooted int       `db:"pancakes_looted" json:"pancakes_looted"`
	DamagePercent  float32   `db:"damage_percent" json:"damage_percent"`
	TimeOfBattle   time.Time `db:"time_of_battle" json:"time_of_battle"`
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
	Width         int
	Height        int
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
	SELECT bd.building_type, bd.building_level, bd.health, COALESCE(dd.building_range, 0), COALESCE(dd.damage_per_sec, 0), COALESCE(dd.damage_per_shot, 0), ob.pos_x, ob.pos_y, bd.width, bd.height
	FROM owned_building ob 
	INNER JOIN building_data bd ON ob.building_data_id = bd.id 
	LEFT JOIN defense_building_data dd ON dd.building_data_id = bd.id 
	WHERE ob.player_id = $1`

	rows, err := database.DB.Query(ctx, query, playerID)
	if err != nil {
		return nil, fmt.Errorf("Error in fetching defense snapshot: %w", err)
	}
	defer rows.Close()

	var defenses []DefenseSnapshot
	for rows.Next() {
		var d DefenseSnapshot
		err = rows.Scan(&d.BuildingType, &d.BuildingLevel, &d.Health, &d.BuildingRange, &d.DamagePerSec, &d.DamagePerShot, &d.PosX, &d.PosY, &d.Width, &d.Height)
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

func UpdateTrophies(ctx context.Context, tx pgx.Tx, newAttackerTrophies int, newDefenderTrophies int, attackerUUID uuid.UUID, defenderUUID uuid.UUID) error {
	query1 := `UPDATE player_info SET trophies = $1 WHERE player_id = $2`
	_, err := tx.Exec(ctx, query1, newAttackerTrophies, attackerUUID)
	if err != nil {
		return fmt.Errorf("Error in updating attacker trophies: %w", err)
	}

	_, err = tx.Exec(ctx, query1, newDefenderTrophies, defenderUUID)
	if err != nil {
		return fmt.Errorf("Error in updating defender trophies: %w", err)
	}

	return nil
}

func UpdateArmyAfterBattle(ctx context.Context, tx pgx.Tx, attackerUUID uuid.UUID) error {
	query2 := `UPDATE trained_troop SET quantity = 0 WHERE player_id = $1`
	_, err := tx.Exec(ctx, query2, attackerUUID)
	if err != nil {
		return fmt.Errorf("Error in clearing attacker's army: %w", err)
	}

	return nil
}

func UpdateShieldTime(ctx context.Context, tx pgx.Tx, shieldEnd time.Time, defenderUUID uuid.UUID) error {
	query3 := `UPDATE player_info SET shield_end_time = $1 WHERE player_id = $2`
	_, err := tx.Exec(ctx, query3, shieldEnd, defenderUUID)
	if err != nil {
		return err
	}

	return nil
}
