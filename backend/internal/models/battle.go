package models

import "time"

type Battle_Log struct {
	ID              string    `json:"id"`
	Attacker_ID     *string   `json:"attacker_id"`
	Defender_ID     *string   `json:"defender_id"`
	Elixir_Looted   int       `json:"elixir_looted"`
	Pancakes_Looted int       `json:"pancakes_looted"`
	Damage_Percent  float32   `json:"damage_percent"`
	Time_Of_Battle  time.Time `json:"time_of_battle"`
}
