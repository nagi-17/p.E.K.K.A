CREATE TABLE trained_troop(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    player_id UUID REFERENCES player_info(player_id) ON DELETE CASCADE,
    troop_data_id int REFERENCES troop_data(id) ON DELETE CASCADE,
    quantity int NOT NULL
);

CREATE TABLE owned_building(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    player_id UUID REFERENCES player_info(player_id) ON DELETE CASCADE,
    building_data_id int REFERENCES building_data(id) ON DELETE CASCADE,
    pos_x int NOT NULL,
    pos_y int NOT NULL,
    upgrade_complete_at TIMESTAMPTZ,
    last_collected_at TIMESTAMPTZ
);

CREATE TABLE player_troop_level(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    player_id UUID REFERENCES player_info(player_id) ON DELETE CASCADE,
    troop_type VARCHAR(50) NOT NULL,
    current_level int NOT NULL DEFAULT 1,
    upgrade_complete_at TIMESTAMPTZ
);
CREATE UNIQUE INDEX idx_player_troop_id ON player_troop_level(player_id, troop_type);