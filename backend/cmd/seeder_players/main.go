package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/nagi-17/p.E.K.K.A/internal/config"
	"github.com/nagi-17/p.E.K.K.A/internal/database"
	"golang.org/x/crypto/bcrypt"
)

type SeedUser struct {
	Username      string
	Email         string
	Password      string
	THLevel       int
	Trophies      int
	SkillPoints   int
	Elixir        int
	Pancakes      int
	TroopLevels   map[string]int
	TrainedTroops map[string]int
}

func main() {
	cfg := config.LoadConfig()
	database.Initialise_DB(cfg.DB_URL)

	ctx := context.Background()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("password123"), 10)
	if err != nil {
		log.Fatalf("Failed to hash password: %v", err)
	}
	hashStr := string(hashedPassword)

	users := []SeedUser{
		{Username: "player_th1_a", Email: "th1_a@pekka.com", Password: "password123", THLevel: 1, Trophies: 100, SkillPoints: 5, Elixir: 1500, Pancakes: 1500,
			TroopLevels:   map[string]int{"Barbarian": 1, "Archer": 1, "Giant": 1, "Goblin": 1, "P.E.K.K.A": 1},
			TrainedTroops: map[string]int{"Barbarian": 5, "Archer": 5}},
		{Username: "player_th1_b", Email: "th1_b@pekka.com", Password: "password123", THLevel: 1, Trophies: 120, SkillPoints: 8, Elixir: 1200, Pancakes: 1200,
			TroopLevels:   map[string]int{"Barbarian": 1, "Archer": 1, "Giant": 1, "Goblin": 1, "P.E.K.K.A": 1},
			TrainedTroops: map[string]int{"Barbarian": 10}},
		{Username: "player_th1_c", Email: "th1_c@pekka.com", Password: "password123", THLevel: 1, Trophies: 150, SkillPoints: 10, Elixir: 1400, Pancakes: 1400,
			TroopLevels:   map[string]int{"Barbarian": 1, "Archer": 1, "Giant": 1, "Goblin": 1, "P.E.K.K.A": 1},
			TrainedTroops: map[string]int{"Archer": 10}},

		{Username: "player_th2_a", Email: "th2_a@pekka.com", Password: "password123", THLevel: 2, Trophies: 420, SkillPoints: 15, Elixir: 5000, Pancakes: 5000,
			TroopLevels:   map[string]int{"Barbarian": 2, "Archer": 2, "Giant": 1, "Goblin": 1, "P.E.K.K.A": 1},
			TrainedTroops: map[string]int{"Barbarian": 10, "Archer": 10}},
		{Username: "player_th2_b", Email: "th2_b@pekka.com", Password: "password123", THLevel: 2, Trophies: 450, SkillPoints: 18, Elixir: 4500, Pancakes: 4500,
			TroopLevels:   map[string]int{"Barbarian": 2, "Archer": 1, "Giant": 1, "Goblin": 1, "P.E.K.K.A": 1},
			TrainedTroops: map[string]int{"Giant": 2, "Barbarian": 10}},
		{Username: "player_th2_c", Email: "th2_c@pekka.com", Password: "password123", THLevel: 2, Trophies: 480, SkillPoints: 20, Elixir: 6000, Pancakes: 6000,
			TroopLevels:   map[string]int{"Barbarian": 1, "Archer": 2, "Giant": 1, "Goblin": 2, "P.E.K.K.A": 1},
			TrainedTroops: map[string]int{"Goblin": 15, "Archer": 5}},

		{Username: "player_th3_a", Email: "th3_a@pekka.com", Password: "password123", THLevel: 3, Trophies: 850, SkillPoints: 40, Elixir: 12000, Pancakes: 12000,
			TroopLevels:   map[string]int{"Barbarian": 3, "Archer": 3, "Giant": 2, "Goblin": 2, "P.E.K.K.A": 1},
			TrainedTroops: map[string]int{"Barbarian": 15, "Archer": 15, "Giant": 2}},
		{Username: "player_th3_b", Email: "th3_b@pekka.com", Password: "password123", THLevel: 3, Trophies: 890, SkillPoints: 45, Elixir: 15000, Pancakes: 15000,
			TroopLevels:   map[string]int{"Barbarian": 3, "Archer": 3, "Giant": 2, "Goblin": 2, "P.E.K.K.A": 1},
			TrainedTroops: map[string]int{"Giant": 4, "Archer": 10}},
		{Username: "player_th3_c", Email: "th3_c@pekka.com", Password: "password123", THLevel: 3, Trophies: 920, SkillPoints: 50, Elixir: 14000, Pancakes: 14000,
			TroopLevels:   map[string]int{"Barbarian": 2, "Archer": 3, "Giant": 2, "Goblin": 3, "P.E.K.K.A": 1},
			TrainedTroops: map[string]int{"Goblin": 20, "Barbarian": 10}},

		{Username: "player_th4_a", Email: "th4_a@pekka.com", Password: "password123", THLevel: 4, Trophies: 1250, SkillPoints: 80, Elixir: 24000, Pancakes: 24000,
			TroopLevels:   map[string]int{"Barbarian": 4, "Archer": 4, "Giant": 3, "Goblin": 3, "P.E.K.K.A": 2},
			TrainedTroops: map[string]int{"Barbarian": 20, "Archer": 20, "Giant": 2}},
		{Username: "player_th4_b", Email: "th4_b@pekka.com", Password: "password123", THLevel: 4, Trophies: 1310, SkillPoints: 90, Elixir: 25000, Pancakes: 25000,
			TroopLevels:   map[string]int{"Barbarian": 4, "Archer": 4, "Giant": 3, "Goblin": 3, "P.E.K.K.A": 2},
			TrainedTroops: map[string]int{"P.E.K.K.A": 1, "Archer": 15}},
		{Username: "player_th4_c", Email: "th4_c@pekka.com", Password: "password123", THLevel: 4, Trophies: 1350, SkillPoints: 100, Elixir: 23000, Pancakes: 23000,
			TroopLevels:   map[string]int{"Barbarian": 3, "Archer": 4, "Giant": 3, "Goblin": 4, "P.E.K.K.A": 2},
			TrainedTroops: map[string]int{"Goblin": 30, "Giant": 2}},
	}

	tx, err := database.DB.Begin(ctx)
	if err != nil {
		log.Fatalf("Failed to begin transaction: %v", err)
	}
	defer tx.Rollback(ctx)

	for _, user := range users {
		log.Printf("Seeding user %s...", user.Username)

		var existingID uuid.UUID
		err = tx.QueryRow(ctx, "SELECT id FROM login_info WHERE username = $1", user.Username).Scan(&existingID)
		if err == nil {
			log.Printf("  User already exists, deleting first...")
			_, err = tx.Exec(ctx, "DELETE FROM login_info WHERE id = $1", existingID)
			if err != nil {
				log.Fatalf("Failed to delete existing user: %v", err)
			}
		}

		var playerID uuid.UUID
		err = tx.QueryRow(ctx, "INSERT INTO login_info (username, email, password_hash) VALUES ($1, $2, $3) RETURNING id",
			user.Username, user.Email, hashStr).Scan(&playerID)
		if err != nil {
			log.Fatalf("Failed to insert login_info: %v", err)
		}

		_, err = tx.Exec(ctx, "INSERT INTO player_info (player_id, trophies, skill_points, elixir, pancakes) VALUES ($1, $2, $3, $4, $5)",
			playerID, user.Trophies, user.SkillPoints, user.Elixir, user.Pancakes)
		if err != nil {
			log.Fatalf("Failed to insert player_info: %v", err)
		}

		for troopType, lvl := range user.TroopLevels {
			_, err = tx.Exec(ctx, "INSERT INTO player_troop_level (player_id, troop_type, current_level, upgrade_complete_at) VALUES ($1, $2, $3, NULL)",
				playerID, troopType, lvl)
			if err != nil {
				log.Fatalf("Failed to insert player_troop_level: %v", err)
			}
		}

		lvl := user.THLevel

		buildingsToPlace := []struct {
			bType string
			bLvl  int
			x, y  int
		}{
			{bType: "Town Hall", bLvl: lvl, x: 18, y: 18},
			{bType: "Cannon", bLvl: lvl, x: 14, y: 14},
			{bType: "Cannon", bLvl: lvl, x: 23, y: 23},

			{bType: "Elixir Collector", bLvl: lvl, x: 23, y: 14},
			{bType: "Pancake Machine", bLvl: lvl, x: 27, y: 18},
			{bType: "Elixir Storage", bLvl: lvl, x: 18, y: 27},
			{bType: "Pancake Stack", bLvl: lvl, x: 10, y: 10},

			{bType: "Army Camp", bLvl: lvl, x: 5, y: 25},
		}

		if lvl >= 2 {
			buildingsToPlace = append(buildingsToPlace, struct {
				bType string
				bLvl  int
				x, y  int
			}{bType: "Archer Tower", bLvl: lvl - 1, x: 14, y: 23})
			buildingsToPlace = append(buildingsToPlace, struct {
				bType string
				bLvl  int
				x, y  int
			}{bType: "Laboratory", bLvl: lvl - 1, x: 25, y: 5})
		}
		if lvl >= 3 {
			buildingsToPlace = append(buildingsToPlace, struct {
				bType string
				bLvl  int
				x, y  int
			}{bType: "Mortar", bLvl: lvl - 2, x: 27, y: 27})
		}

		for _, b := range buildingsToPlace {

			var bDataID int
			err = tx.QueryRow(ctx, "SELECT id FROM building_data WHERE building_type = $1 AND building_level = $2", b.bType, b.bLvl).Scan(&bDataID)
			if err != nil {
				log.Fatalf("Could not find building_data ID for %s lvl %d: %v", b.bType, b.bLvl, err)
			}

			_, err = tx.Exec(ctx, "INSERT INTO owned_building (player_id, building_data_id, pos_x, pos_y, upgrade_complete_at, last_collected_at) VALUES ($1, $2, $3, $4, NULL, $5)",
				playerID, bDataID, b.x, b.y, time.Now())
			if err != nil {
				log.Fatalf("Failed to insert owned_building: %v", err)
			}
		}

		for troopType, qty := range user.TrainedTroops {
			userLvl := user.TroopLevels[troopType]
			var tDataID int
			err = tx.QueryRow(ctx, "SELECT id FROM troop_data WHERE troop_type = $1 AND troop_level = $2", troopType, userLvl).Scan(&tDataID)
			if err != nil {
				log.Fatalf("Could not find troop_data ID for %s lvl %d: %v", troopType, userLvl, err)
			}

			_, err = tx.Exec(ctx, "INSERT INTO trained_troop (player_id, troop_data_id, quantity) VALUES ($1, $2, $3)",
				playerID, tDataID, qty)
			if err != nil {
				log.Fatalf("Failed to insert trained_troop: %v", err)
			}
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		log.Fatalf("Failed to commit transaction: %v", err)
	}

	fmt.Println("SUCCESS: Successfully seeded 12 dummy users into the database!")
}
