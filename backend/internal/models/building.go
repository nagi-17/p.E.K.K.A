package models

import "time"

type Building_Data struct {
	ID                     int           `json:"id"`
	Building_Type          string        `json:"building_type"`
	Building_Level         int           `json:"building_level"`
	Health                 int           `json:"health"`
	Width                  int           `json:"width"`
	Height                 int           `json:"height"`
	Build_Time             time.Duration `json:"build_time"`
	Upgrade_Cost_Elixir    int           `json:"upgrade_cost_elixir"`
	Upgrade_Cost_Pancakes  int           `json:"upgrade_cost_pancakes"`
	Upgrade_Time           time.Duration `json:"upgrade_time"`
	Max_Quantity_Available int           `json:"max_quantity_available"`
	Skill_On_Upgrade       int           `json:"skill_on_upgrade"`
}

type Owned_Building_Data struct {
	ID                  string     `json:"id"`
	Player_ID           string     `json:"player_id"`
	Building_Data_ID    int        `json:"building_data_id"`
	Pos_X               int        `json:"pos_x"`
	Pos_Y               int        `json:"pos_y"`
	Upgrade_Complete_At *time.Time `json:"upgrade_complete_at"`
	Last_Collected_At   *time.Time `json:"last_collected_at"`
}

type Town_Hall_Data struct {
	Building_Data_ID                int `json:"building_data_id"`
	Min_Skill_Points_Before_Upgrade int `json:"min_skill_points_before_upgrade"`
}

type Defense_Building_Data struct {
	Building_Data_ID           int `json:"building_data_id"`
	Range                      int `json:"range"`
	Damage                     int `json:"damage"`
	Max_Possible_Upgrade_Level int `json:"max_possible_upgrade_level"`
	Unlock_Town_Hall_Level     int `json:"unlock_town_hall_level"`
}

type Resource_Building_Data struct {
	Building_Data_ID           int `json:"building_data_id"`
	Elixir_Gen_Per_Min         int `json:"elixir_gen_per_min"`
	Pancakes_Gen_Per_Min       int `json:"pancakes_gen_per_min"`
	Max_Possible_Upgrade_Level int `json:"max_possible_upgrade_level"`
	Unlock_Town_Hall_Level     int `json:"unlock_town_hall_level"`
}

type Storage_Building_Data struct {
	Building_Data_ID           int `json:"building_data_id"`
	Max_Storage                int `json:"max_storage"`
	Max_Possible_Upgrade_Level int `json:"max_possible_upgrade_level"`
	Unlock_Town_Hall_Level     int `json:"unlock_town_hall_level"`
}

type Laboratory_Data struct {
	Building_Data_ID           int `json:"building_data_id"`
	Max_Troop_Upgrade_Level    int `json:"max_troop_upgrade_level"`
	Max_Possible_Upgrade_Level int `json:"max_possible_upgrade_level"`
	Unlock_Town_Hall_Level     int `json:"unlock_town_hall_level"`
}
