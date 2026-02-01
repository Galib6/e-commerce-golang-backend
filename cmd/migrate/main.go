package main

import (
	"log"
	"path/filepath"

	"github.com/goutamkumar/golang_restapi_postgresql_test1/internal/config"
)

func main() {
	// 1. Load config
	cfg := config.LoadEnv()

	// Default migrations dir if not set
	if cfg.Migrations == "" {
		cfg.Migrations = "migrations"
	}

	// 2. Connect to Database
	db, err := config.Connect(cfg.DatabaseUrl)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Resolve absolute path for migrations
	absPath, err := filepath.Abs(cfg.Migrations)
	if err != nil {
		log.Fatalf("Failed to resolve migrations path: %v", err)
	}

	// 3. Run Migrations
	log.Printf("Running migrations from: %s", absPath)
	if err := config.RunMigrations(db, absPath); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	log.Println("✅ Migrations completed successfully")
}
