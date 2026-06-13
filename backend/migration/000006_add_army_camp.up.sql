CREATE TABLE amry_camp_data {
    building_data_id int PRIMARY KEY REFERENCES building_data(id) ON DELETE CASCADE,
    housing_space int NOT NULL,
    max_possible_upgrade_level int NOT NULL,
    unlock_town_hall_level INT DEFAULT 1
}