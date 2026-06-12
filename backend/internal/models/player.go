package models

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/nagi-17/p.E.K.K.A/internal/database"
)

type LoginInfo struct {
	ID            uuid.UUID `db:"id"`
	Username      string    `db:"username"`
	Email         string    `db:"email"`
	Password_Hash string    `db:"password_hash"`
	Created_At    time.Time `db:"created_at"`
}

type PlayerInfo struct {
	Player_ID       uuid.UUID  `db:"player_id"`
	Trophies        int        `db:"trophies"`
	Skill_points    int        `db:"skill_points"`
	Elixir          int        `db:"elixir"`
	Pancakes        int        `db:"pancakes"`
	Shield_End_Time *time.Time `db:"shield_end_time"`
}

func RegisterNewPlayer(ctx context.Context, username string, email string, pass_hash string) (string, error) {
	tx, err := database.DB.Begin(ctx)
	if err != nil {
		return "", fmt.Errorf("Failed to begin databse transaction: %w", err)
	}

	defer tx.Rollback(ctx)

	login_info_query := `
		INSERT INTO login_info (username, email, password_hash)
		VALUES ($1, $2, $3)
		RETURNING id;
	`

	var player_ID uuid.UUID
	err = tx.QueryRow(ctx, login_info_query, username, email, pass_hash).Scan(&player_ID)
	if err != nil {
		return "", fmt.Errorf("Failed to register user: %w", err)
	}

	player_info_query := `
		INSERT INTO player_info (player_id, trophies, skill_points, elixir, pancakes, shield_end_time)
		VALUES ($1, $2, $3, $4, $5, $6);
	`

	var init_trophies int = 0
	var init_skill_points int = 0
	var init_elixir int = 1000
	var init_pancakes int = 1000
	var init_shield_time *time.Time = nil
	_, err = tx.Exec(ctx, player_info_query, player_ID, init_trophies, init_skill_points, init_elixir, init_pancakes, init_shield_time)
	if err != nil {
		return "", fmt.Errorf("Error in loading initial player data: %w", err)
	}

	var thDataID int
	query := `SELECT id FROM building_data WHERE building_type = 'Town Hall' AND building_level = 1`

	err = tx.QueryRow(ctx, query).Scan(&thDataID)
	if err != nil {
		return "", fmt.Errorf("Cannot find Level 1 Town Hall in db: %w", err)
	}

	insertTHquery := `INSERT INTO owned_building (id, player_id, building_data_id, pos_x, pos_y, upgrade_complete_at, last_collected_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7)`
	thID := uuid.New()
	var instantBuildth *time.Time = nil
	var noCollect *time.Time = nil

	_, err = tx.Exec(ctx, insertTHquery, thID, player_ID, thDataID, 20, 20, instantBuildth, noCollect)
	if err != nil {
		return "", fmt.Errorf("Failed to place default town hall: %w", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return "", fmt.Errorf("Failed to commit to database: %w", err)
	}
	return player_ID.String(), nil
}

func GetLoginInfoUsingUsername(ctx context.Context, username string) (*LoginInfo, error) {
	query := `
		SELECT id, username, email, password_hash, created_at 
		FROM login_info
		WHERE username = $1
		`

	var info LoginInfo
	err := database.DB.QueryRow(ctx, query, username).Scan(&info.ID, &info.Username, &info.Email, &info.Password_Hash, &info.Created_At)
	if err != nil {
		return nil, fmt.Errorf("Error in fetching user data: %w", err)
	}

	return &info, nil
}

func GetPlayerInfoByID(ctx context.Context, playerID uuid.UUID) (*PlayerInfo, error) {
	query := `
	SELECT player_id, trophies, skill_points, elixir, pancakes, shield_end_time
	FROM player_info
	WHERE player_id=$1
	`
	var info PlayerInfo
	err := database.DB.QueryRow(ctx, query, playerID).Scan(&info.Player_ID, &info.Trophies, &info.Skill_points, &info.Elixir, &info.Pancakes, &info.Shield_End_Time)
	if err != nil {
		return nil, fmt.Errorf("Error in fetching user stats: %w", err)
	}
	return &info, nil
}

func GetPlayerTownHallLevel(ctx context.Context, playerID uuid.UUID) (int, error) {
	query := `
	SELECT bd.building_level
	FROM owned_building ob
	INNER JOIN building_data bd ON ob.building_data_id = bd.id
	WHERE ob.player_id = $1 AND bd.building_type = 'Town Hall'
	LIMIT 1
	`
	var level int
	err := database.DB.QueryRow(ctx, query, playerID).Scan(&level)
	if err != nil {
		return 0, fmt.Errorf("Failed to fetch Town Hall level")
	}
	return level, nil
}

func UpdateResources(ctx context.Context, tx pgx.Tx, playerID uuid.UUID, newPancakes int, newElixir int) error {

	maxElixir, maxPancakes, err := GetPlayerStorageCapacity(ctx, playerID)
	if err != nil {
		return err
	}
	if newElixir > maxElixir {
		newElixir = maxElixir
	}
	if newPancakes > maxPancakes {
		newPancakes = maxPancakes
	}

	query := `
	UPDATE player_info
	SET elixir = $1, pancakes = $2
	WHERE player_id = $3
	`
	_, err = tx.Exec(ctx, query, newElixir, newPancakes, playerID)
	if err != nil {
		return fmt.Errorf("Error in updating resources")
	}
	return nil
}

func AddSkill(ctx context.Context, tx pgx.Tx, playerID uuid.UUID, newSkill int) error {

	query := `UPDATE player_info SET skill_points = $1 WHERE player_id = $2`
	_, err := tx.Exec(ctx, query, newSkill, playerID)

	if err != nil {
		return fmt.Errorf("Error in adding skill: %w", err)
	}

	return nil
}
