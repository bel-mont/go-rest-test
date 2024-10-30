package db

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"log"
	"os"
	"time"
)

// InitDB initializes a pgx connection pool for PostgreSQL.
func InitDB() *pgxpool.Pool {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found, using default environment variables")
	}

	// Retrieve database connection details from environment variables
	dbName := os.Getenv("DB_NAME")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")

	// Format connection string
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", dbUser, dbPassword, dbHost, dbPort, dbName)

	// Configure pool settings if needed
	config, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		log.Fatal("Unable to parse database configuration:", err)
	}
	//config.MaxConns = 10
	config.HealthCheckPeriod = 1 * time.Minute

	// Establish the connection pool
	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		log.Fatal("Failed to create connection pool:", err)
	}

	// Check if the connection is working
	if err := pool.Ping(context.Background()); err != nil {
		log.Fatal("Database connection is not active:", err)
	}

	return pool
}
