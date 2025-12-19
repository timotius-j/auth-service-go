package config

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"
)

func RunMigrations(db *sql.DB) {
	if os.Getenv("APP_ENV") != "local" {
		log.Println("⚠️ skipping auto migrations (not local env)")
		return
	}

	files := []string{
		"migrations/000_drop_tables.sql",
		"migrations/001_schema.sql",
		"migrations/002_seeder.sql",
	}

	for _, file := range files {
		content, err := os.ReadFile(filepath.Clean(file))
		if err != nil {
			log.Fatalf("failed to read migration %s: %v", file, err)
		}

		if _, err := db.Exec(string(content)); err != nil {
			log.Fatalf("failed to run migration %s: %v", file, err)
		}

		log.Printf("✅ executed migration: %s", file)
	}
}
