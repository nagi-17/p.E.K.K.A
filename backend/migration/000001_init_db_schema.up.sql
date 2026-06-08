CREATE TABLE login_info(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT now()
);

CREATE TABLE player_info(
    player_id UUID PRIMARY KEY REFERENCES login_info(id) ON DELETE CASCADE,
    trophies int NOT NULL DEFAULT 0,
    skill_points int NOT NULL DEFAULT 0,
    elixir int NOT NULL DEFAULT 0,
    pancakes int NOT NULL DEFAULT 0,
    shield_end_time TIMESTAMPTZ
);
