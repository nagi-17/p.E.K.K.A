package models

import "time"

type Login_info struct {
	ID            string    `json:"id"`
	Username      string    `json:"username"`
	Email         string    `json:"email"`
	Password_Hash string    `json:"-"`
	Created_At    time.Time `json:"created_at"`
}

type Player_info struct {
	Player_ID       string     `json:"player_id"`
	Trophies        int        `json:"trophies"`
	Skill_points    int        `json:"skill_points"`
	Elixir          int        `json:"elixir"`
	Pancakes        int        `json:"pancakes"`
	Shield_End_Time *time.Time `json:"sheild_end_at"`
}
