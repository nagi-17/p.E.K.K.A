CREATE TABLE troop_data(
    id int PRIMARY KEY,
    troop_type VARCHAR(50) NOT NULL,
    troop_level int DEFAULT 1 NOT NULL,
    health int NOT NULL,
    damage int NOT NULL,
    troop_range int NOT NULL,
    training_cost_elixir int NOT NULL,
    space_occupied_in_army int NOT NULL,
    speed int NOT NULL
);
CREATE UNIQUE INDEX idx_troop_data_type_level ON troop_data (troop_type, troop_level);

CREATE TABLE building_data(
    id int PRIMARY KEY,
    building_type VARCHAR(50) NOT NULL,
    building_level int DEFAULT 1,
    health int,
    width int NOT NULL,
    height int NOT NULL,
    build_time int NOT NULL,
    upgrade_cost_elixir int DEFAULT 0,
    upgrade_cost_pancakes int DEFAULT 0,
    upgrade_time int NOT NULL,
    max_quantity_available int NOT NULL
);
CREATE UNIQUE INDEX idx_building_data_type_level ON building_data (building_type, building_level);

CREATE TABLE town_hall_data(
    building_data_id int PRIMARY KEY REFERENCES building_data(id) ON DELETE CASCADE,
    min_skill_points_before_upgrade int NOT NULL
);

CREATE TABLE defense_building_data(
    building_data_id int PRIMARY KEY REFERENCES building_data(id) ON DELETE CASCADE,
    building_range int NOT NULL,
    damage int NOT NULL,
    max_possible_upgrade_level int NOT NULL,
    unlock_town_hall_level int NOT NULL
);

CREATE TABLE resource_building_data(
    building_data_id int PRIMARY KEY REFERENCES building_data(id) ON DELETE CASCADE,
    elixir_gen_per_min int NOT NULL,
    pancakes_gen_per_min int NOT NULL,
    max_possible_upgrade_level int NOT NULL,
    unlock_town_hall_level int DEFAULT 1
);

CREATE TABLE storage_building_data(
    building_data_id int PRIMARY KEY REFERENCES building_data(id) ON DELETE CASCADE,
    max_storage int NOT NULL,
    max_possible_upgrade_level int NOT NULL,
    unlock_town_hall_level int DEFAULT 1
);

CREATE TABLE laboratory_data(
    building_data_id int PRIMARY KEY REFERENCES building_data(id) ON DELETE CASCADE,
    max_troop_upgrade_level int NOT NULL,
    unlock_town_hall_level int NOT NULL,
    max_possible_upgrade_level int NOT NULL
);