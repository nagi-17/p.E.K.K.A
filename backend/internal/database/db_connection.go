package database

import (
	"context"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

var DB *pgxpool.Pool

func Initialise_DB(db_url string) {
	ctx := context.Background()

	config, err := pgxpool.ParseConfig(db_url)
	if err != nil {
		log.Fatalf("Unable to parse connection string: %v\n", err)
	}
	config.MaxConns = 20
	config.MinConns = 5
	config.MaxConnLifetime = 5 * time.Minute
	config.MaxConnIdleTime = 2 * time.Minute
	config.HealthCheckPeriod = 1 * time.Minute

	DB, err = pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		log.Fatalf("Unable to create connection pool: %v\n", err)
	}
	err = DB.Ping(ctx)
	if err != nil {
		log.Fatalf("Unable to ping the database: %v\n", err)
	}
	log.Println("Successfully connected to database")
}
