package db

import (
	"database/sql"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func Connect() (*sql.DB, error) {
	// Load .env
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		log.Fatal("DATABASE_URL not set")
	}

	database, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	// Set default schema to public
	if _, err := database.Exec("SET search_path TO public;"); err != nil {
		return nil, err
	}

	// Test connection
	if err := database.Ping(); err != nil {
		return nil, err
	}

	return database, nil
}

