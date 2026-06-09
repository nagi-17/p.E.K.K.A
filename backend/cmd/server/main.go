package main

import (
	"log"
	"net/http"

	"github.com/nagi-17/p.E.K.K.A/internal/config"
	"github.com/nagi-17/p.E.K.K.A/internal/database"
	"github.com/nagi-17/p.E.K.K.A/internal/routes"
)

func main() {
	config_data := config.LoadConfig()

	database.Initialise_DB(config_data.DB_URL)

	log.Printf("Server starting on port: %s\n", config_data.Port)
	router := routes.InitRouter()

	err := http.ListenAndServe(":"+config_data.Port, router)
	if err != nil {
		log.Fatalf("Server has crashed: %v\n", err)
	}
}
