CREATE TABLE battle_log(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    attacker_id UUID REFERENCES player_info(player_id) ON DELETE SET NULL,
    defender_id UUID REFERENCES player_info(player_id) ON DELETE SET NULL,
    elixir_looted int NOT NULL,
    pancakes_looted int NOT NULL,
    trophy_change_attacker int NOT NULL,
    trophy_change_defender int NOT NULL,
    damage_percent float NOT NULL,
    time_of_battle TIMESTAMPTZ DEFAULT now()
);