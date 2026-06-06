package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port      string
	DBHost    string
	JWTSecret string
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}

	server_port := os.Getenv("SERVER_PORT")
	db_host := os.Getenv("DB_HOST")
	jwt_secret := os.Getenv("JWT_SECRET")

	return &Config{
		Port:      server_port,
		DBHost:    db_host,
		JWTSecret: jwt_secret,
	}
}
