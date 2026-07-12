package db

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func DBConfig() *pgxpool.Pool {
	// Initialize Postgres DB if env vars are present
	pgUser := os.Getenv("POSTGRES_USER")
	pgPass := os.Getenv("POSTGRES_PASSWORD")
	pgDB := os.Getenv("POSTGRES_DB")
	pgHost := os.Getenv("POSTGRES_HOST")
	pgPort := os.Getenv("POSTGRES_PORT")

	var pool *pgxpool.Pool

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", pgUser, pgPass, pgHost, pgPort, pgDB)

	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		log.Panicf("failed to parse connection config: %v", err)
	}

	config.MaxConns = 50
	config.MinConns = 10
	config.MaxConnLifetime = time.Hour
	config.MaxConnIdleTime = time.Minute
	config.HealthCheckPeriod = time.Minute

	pool, err = pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		log.Panicf("failed to create connection pool: %v", err)
	}
	if err := pool.Ping(context.Background()); err != nil {
		log.Panicf("failed to ping db: %v", err)
	}
	fmt.Printf("connected to Postgres at %s:%s/%s", pgHost, pgPort, pgDB)

	// Middleware to attach db pool to context
	return pool

}
