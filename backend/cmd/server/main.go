package main

import (
	"log"

	"github.com/nagi-17/p.E.K.K.A/internal/config"
	"github.com/nagi-17/p.E.K.K.A/internal/database"
)

func main() {
	config_data := config.LoadConfig()

	database.Initialise_DB(config_data.DB_URL)

	log.Printf("Server starting on port: %s\n", config_data.Port)
}
