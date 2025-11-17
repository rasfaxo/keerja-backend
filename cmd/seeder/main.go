package main

import (
	"log"
	"os"

	"keerja-backend/database/seeders"
	"keerja-backend/internal/config"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Initialize database
	db, err := config.InitDB(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Run seeders
	if err := seeders.RunSeeders(db); err != nil {
		log.Fatalf("Failed to run seeders: %v", err)
		os.Exit(1)
	}

	log.Println("Seeders completed successfully!")
}
