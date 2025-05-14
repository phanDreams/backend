package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq" // PostgreSQL driver
)

var DB *sql.DB

func InitDB() error {
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		return fmt.Errorf("DATABASE_URL environment variable is not set")
	}

	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("error opening database: %v", err)
	}

	// Set connection pool settings
	DB.SetMaxOpenConns(25)
	DB.SetMaxIdleConns(25) 

	// Verify the connection
	err = DB.Ping()
	if err != nil {
		return fmt.Errorf("error connecting to the database: %v", err)
	}

	log.Println("✅ Successfully connected to Supabase Postgres!")
	return nil
}

// GetDB returns the database instance
func GetDB() *sql.DB {
	return DB
} 