package models

import "time"

type Troop_Data struct {
	ID                     int    `json:"id"`
	Troop_Type             string `json:"troop_type"`
	Troop_Level            int    `json:"troop_level"`
	Health                 int    `json:"health"`
	Damage                 int    `json:"damage"`
	Troop_Range            int    `json:"troop_range"`
	Training_Cost_Elixir   int    `json:"training_cost_elixir"`
	Space_Occupied_In_Army int    `json:"space_occupied_in_army"`
	Speed                  int    `json:"speed"`
}

type Trained_Troop struct {
	ID            string `json:"id"`
	Player_ID     string `json:"player_id"`
	Troop_Data_ID int    `json:"troop_data_id"`
	Quantity      int    `json:"quantity"`
}

type Player_Troop_Level struct {
	ID                  string    `json:"id"`
	Player_ID           string    `json:"player_id"`
	Troop_Type          string    `json:"troop_type"`
	Current_Level       int       `json:"current_level"`
	Upgrade_Complete_At time.Time `json:"upgrade_complete_at"`
}
