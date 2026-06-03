CREATE TABLE login_info(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT now()
);

CREATE TABLE player_info(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    player_id UUID UNIQUE   REFERENCES login_info(id) ON DELETE CASCADE,
    player_tag VARCHAR(20) UNIQUE NOT NULL,
    trophies int DEFAULT 0,
    skill_points int DEFAULT 0,
    elixir int DEFAULT 0,
    pancakes int DEFAULT 0,
    shield_end_time TIMESTAMPTZ
);

CREATE TABLE troop_data(
    id int PRIMARY KEY,
    troop_type VARCHAR(50) NOT NULL,
    troop_level int DEFAULT 1 NOT NULL,
    health int NOT NULL,
    damage int NOT NULL,
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
    damage int DEFAULT 0,
    width int NOT NULL,
    height int NOT NULL,
    build_time int NOT NULL,
    cost_elixir int DEFAULT 0,
    cost_pancakes int DEFAULT 0
);
CREATE UNIQUE INDEX idx_building_data_type_level ON building_data (building_type, building_level);

CREATE TABLE trained_troop(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    player_id UUID REFERENCES player_info(id) ON DELETE CASCADE,
    troop_data_id int REFERENCES troop_data(id) ON DELETE CASCADE,
    quantity int NOT NULL
);

CREATE TABLE owned_building(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    player_id UUID REFERENCES player_info(id) ON DELETE CASCADE,
    building_data_id int REFERENCES building_data(id) ON DELETE CASCADE,
    pos_x int NOT NULL,
    pos_y int NOT NULL,
    upgrade_time TIMESTAMPTZ
);

CREATE TABLE battle_log(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    attacker_id UUID REFERENCES player_info(id) ON DELETE SET NULL,
    defender_id UUID REFERENCES player_info(id) ON DELETE SET NULL,
    elixir_looted int NOT NULL,
    pancakes_looted int NOT NULL,
    trophy_change_attacker int NOT NULL,
    trophy_change_defender int NOT NULL,
    damage_percent int NOT NULL,
    time_of_battle TIMESTAMPTZ DEFAULT now()
);