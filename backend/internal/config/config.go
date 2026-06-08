package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port      string
	DB_URL    string
	JWTSecret string
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}

	server_port := os.Getenv("SERVER_PORT")
	db_url := os.Getenv("DB_URL")
	jwt_secret := os.Getenv("JWT_SECRET")

	return &Config{
		Port:      server_port,
		DB_URL:    db_url,
		JWTSecret: jwt_secret,
	}
}
