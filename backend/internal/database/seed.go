package database

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

func SeedDatabase(pool *pgxpool.Pool) {
	ctx := context.Background()

	var count int
	err := pool.QueryRow(ctx, "SELECT COUNT(*) FROM troop_data").Scan(&count)
	if err == nil && count > 0 {
		log.Println("Static game data already exists. Skipping seeder.")
		return
	}

	insertTroops := `
		INSERT INTO troop_data (id, troop_type, troop_level, health, dps, dpa, troop_range, upgrade_cost_elixir, space_occupied_in_army, mov_speed, attack_speed, lab_level_req, airborne, upgrade_time) VALUES

		(101, 'Barbarian', 1, 45, 8, 8, 1, 0, 1, 16, 1, 1, FALSE, 0),
		(102, 'Barbarian', 2, 54, 11, 11, 1, 50000, 1, 16, 1, 1, FALSE, 60),
		(103, 'Barbarian', 3, 65, 14, 14, 1, 150000, 1, 16, 1, 2, FALSE, 300),
		(104, 'Barbarian', 4, 78, 18, 18, 1, 500000, 1, 16, 1, 3, FALSE, 900),
		
		(201, 'Archer', 1, 30, 7, 7, 3, 0, 1, 24, 1, 1, FALSE, 0),
		(202, 'Archer', 2, 36, 9, 9, 3, 50000, 1, 24, 1, 1, FALSE, 60),
		(203, 'Archer', 3, 44, 12, 12, 3, 250000, 1, 24, 1, 2, FALSE, 600),
		(204, 'Archer', 4, 54, 16, 16, 3, 750000, 1, 24, 1, 3, FALSE, 1200),

		(301, 'Giant', 1, 300, 11, 22, 1, 0, 5, 12, 2, 1, FALSE, 0),
		(302, 'Giant', 2, 360, 14, 28, 1, 100000, 5, 12, 2, 1, FALSE, 300),
		(303, 'Giant', 3, 430, 19, 38, 1, 250000, 5, 12, 2, 2, FALSE, 900),
		(304, 'Giant', 4, 520, 24, 48, 1, 750000, 5, 12, 2, 3, FALSE, 1800),

		(401, 'Goblin', 1, 25, 11, 11, 1, 0, 1, 32, 1, 1, FALSE, 0),
		(402, 'Goblin', 2, 30, 14, 14, 1, 50000, 1, 32, 1, 1, FALSE, 60),
		(403, 'Goblin', 3, 36, 19, 19, 1, 250000, 1, 32, 1, 2, FALSE, 600),
		(404, 'Goblin', 4, 43, 24, 24, 1, 750000, 1, 32, 1, 3, FALSE, 1200),

		(501, 'P.E.K.K.A', 1, 2800, 240, 600, 1, 0, 25, 16, 2, 1, FALSE, 0),
		(502, 'P.E.K.K.A', 2, 3100, 270, 675, 1, 3000000, 25, 16, 2, 2, FALSE, 600),
		(503, 'P.E.K.K.A', 3, 3500, 310, 775, 1, 6000000, 25, 16, 2, 3, FALSE, 1800),
		(504, 'P.E.K.K.A', 4, 4000, 360, 900, 1, 8000000, 25, 16, 2, 4, FALSE, 3600);
	`
	if _, err = pool.Exec(ctx, insertTroops); err != nil {
		log.Printf("Error seeding troop data: %v", err)
		return
	}

	insertBuildings := `
		INSERT INTO building_data (id, building_type, building_level, health, width, height, build_time, upgrade_cost_elixir, upgrade_cost_pancakes, upgrade_time, max_quantity_available, skill_on_upgrade) VALUES

		(1001, 'Town Hall', 1, 1500, 4, 4, 0, 0, 0, 0, 1, 0),
		(1002, 'Town Hall', 2, 1600, 4, 4, 600, 0, 1000, 600, 1, 10),
		(1003, 'Town Hall', 3, 1850, 4, 4, 3600, 0, 4000, 3600, 1, 25),
		(1004, 'Town Hall', 4, 2100, 4, 4, 14400, 0, 25000, 14400, 1, 50),

		(2001, 'Cannon', 1, 400, 3, 3, 60, 0, 250, 60, 2, 5),
		(2002, 'Cannon', 2, 450, 3, 3, 900, 0, 1000, 900, 2, 10),
		(2003, 'Cannon', 3, 500, 3, 3, 2700, 0, 4000, 2700, 2, 15),
		(2004, 'Cannon', 4, 570, 3, 3, 14400, 0, 16000, 14400, 2, 25),

		(2101, 'Archer Tower', 1, 380, 3, 3, 900, 0, 1000, 900, 1, 10),
		(2102, 'Archer Tower', 2, 420, 3, 3, 2700, 0, 2000, 2700, 1, 15),
		(2103, 'Archer Tower', 3, 460, 3, 3, 14400, 0, 5000, 14400, 1, 25),
		(2104, 'Archer Tower', 4, 500, 3, 3, 43200, 0, 20000, 43200, 1, 50),

		(2201, 'Mortar', 1, 400, 3, 3, 14400, 0, 8000, 14400, 1, 25),
		(2202, 'Mortar', 2, 450, 3, 3, 43200, 0, 32000, 43200, 1, 50),
		(2203, 'Mortar', 3, 500, 3, 3, 86400, 0, 120000, 86400, 1, 100),
		(2204, 'Mortar', 4, 550, 3, 3, 172800, 0, 400000, 172800, 1, 200),

		(3001, 'Elixir Collector', 1, 400, 3, 3, 60, 0, 150, 60, 2, 5),
		(3002, 'Elixir Collector', 2, 450, 3, 3, 600, 0, 300, 600, 2, 10),
		(3003, 'Elixir Collector', 3, 500, 3, 3, 2700, 0, 700, 2700, 2, 15),
		(3004, 'Elixir Collector', 4, 550, 3, 3, 14400, 0, 1400, 14400, 2, 25),

		(3101, 'Pancake Machine', 1, 400, 3, 3, 60, 150, 0, 60, 2, 5),
		(3102, 'Pancake Machine', 2, 450, 3, 3, 600, 300, 0, 600, 2, 10),
		(3103, 'Pancake Machine', 3, 500, 3, 3, 2700, 700, 0, 2700, 2, 15),
		(3104, 'Pancake Machine', 4, 550, 3, 3, 14400, 1400, 0, 14400, 2, 25),

		(4001, 'Elixir Storage', 1, 400, 3, 3, 60, 0, 300, 60, 1, 5),
		(4002, 'Elixir Storage', 2, 600, 3, 3, 1800, 0, 750, 1800, 1, 10),
		(4003, 'Elixir Storage', 3, 800, 3, 3, 7200, 0, 1500, 7200, 1, 15),
		(4004, 'Elixir Storage', 4, 1000, 3, 3, 14400, 0, 3000, 14400, 1, 25),

		(4101, 'Pancake Stack', 1, 400, 3, 3, 60, 300, 0, 60, 1, 5),
		(4102, 'Pancake Stack', 2, 600, 3, 3, 1800, 750, 0, 1800, 1, 10),
		(4103, 'Pancake Stack', 3, 800, 3, 3, 7200, 1500, 0, 7200, 1, 15),
		(4104, 'Pancake Stack', 4, 1000, 3, 3, 14400, 3000, 0, 14400, 1, 25),

		(5001, 'Laboratory', 1, 500, 3, 3, 1800, 25000, 0, 1800, 1, 15),
		(5002, 'Laboratory', 2, 550, 3, 3, 14400, 50000, 0, 14400, 1, 25),
		(5003, 'Laboratory', 3, 600, 3, 3, 43200, 90000, 0, 43200, 1, 50),
		(5004, 'Laboratory', 4, 650, 3, 3, 86400, 250000, 0, 86400, 1, 100),

		(6001, 'Army Camp', 1, 400, 5, 5, 300, 250, 0, 300, 1, 5),
		(6002, 'Army Camp', 2, 500, 5, 5, 1800, 2500, 0, 1800, 1, 10),
		(6003, 'Army Camp', 3, 600, 5, 5, 7200, 10000, 0, 7200, 1, 15),
		(6004, 'Army Camp', 4, 700, 5, 5, 28800, 100000, 0, 28800, 1, 25);
	`
	if _, err = pool.Exec(ctx, insertBuildings); err != nil {
		log.Printf("Error seeding building data: %v", err)
		return
	}

	insertSubTables := `
		INSERT INTO town_hall_data (building_data_id, min_skill_points_before_upgrade) VALUES
		(1001, 0), (1002, 50), (1003, 150), (1004, 300);

		INSERT INTO defense_building_data (building_data_id, building_range, damage_per_sec, damage_per_shot, max_possible_upgrade_level, unlock_town_hall_level) VALUES
		(2001, 9, 9, 7, 4, 1), (2002, 9, 11, 9, 4, 1), (2003, 9, 15, 12, 4, 2), (2004, 9, 19, 15, 4, 3),
		(2101, 10, 11, 11, 4, 2), (2102, 10, 15, 15, 4, 2), (2103, 10, 19, 19, 4, 3), (2104, 10, 25, 25, 4, 4),
		(2201, 11, 4, 20, 4, 3), (2202, 11, 5, 25, 4, 3), (2203, 11, 6, 30, 4, 4), (2204, 11, 8, 40, 4, 4);

		INSERT INTO resource_building_data (building_data_id, elixir_gen_per_min, pancakes_gen_per_min, max_possible_upgrade_level, unlock_town_hall_level) VALUES
		(3001, 20, 0, 4, 1), (3002, 40, 0, 4, 1), (3003, 60, 0, 4, 2), (3004, 80, 0, 4, 2),
		(3101, 0, 20, 4, 1), (3102, 0, 40, 4, 1), (3103, 0, 60, 4, 2), (3104, 0, 80, 4, 2);

		INSERT INTO storage_building_data (building_data_id, max_storage, max_possible_upgrade_level, unlock_town_hall_level) VALUES
		(4001, 1500, 4, 1), (4002, 3000, 4, 1), (4003, 10000, 4, 2), (4004, 25000, 4, 3),
		(4101, 1500, 4, 1), (4102, 3000, 4, 1), (4103, 10000, 4, 2), (4104, 25000, 4, 3);

		INSERT INTO laboratory_data (building_data_id, max_troop_upgrade_level, unlock_town_hall_level, max_possible_upgrade_level) VALUES
		(5001, 1, 2, 4), (5002, 2, 3, 4), (5003, 3, 4, 4), (5004, 4, 4, 4);

		INSERT INTO army_camp_data (building_data_id, housing_space, max_possible_upgrade_level, unlock_town_hall_level) VALUES
		(6001, 20, 4, 1), (6002, 30, 4, 2), (6003, 40, 4, 3), (6004, 50, 4, 4);
	`
	if _, err = pool.Exec(ctx, insertSubTables); err != nil {
		log.Printf("Error seeding sub-tables: %v", err)
		return
	}

	log.Println("Database successfully seeded with static game data")
}
