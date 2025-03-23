package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq" // Import the driver (PostgreSQL example)
	"github.com/pressly/goose/v3"
)

func main() {
	// Replace with your database connection string
	db, err := sql.Open("postgres", "postgres://test:test@postgres:5432/testdb?sslmode=disable")
	if err != nil {
		log.Fatal("Failed to connect to DB:", err)
	}
	defer db.Close()

	// Set the path to the migrations folder
	migrationsDir := "migrations"

	// Run all up migrations
	if err := goose.Up(db, migrationsDir); err != nil {
		log.Printf("Migration failed: %v", err)
	}

	log.Println("Migrations applied successfully!")
}
