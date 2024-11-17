package db

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"io"
	"log"
	"os"
	"path/filepath"
)

// ApplyMigrations applies all SQL migrations from the migrations folder to the provided DB connection.
func ApplyMigrations(dbConn *sqlx.DB) error {
	migrationsDir := "./migrations"
	dir, err := os.Open(migrationsDir)
	if err != nil {
		return fmt.Errorf("failed to open migrations directory: %w", err)
	}
	defer dir.Close()

	// Read all entries in the migrations directory
	files, err := dir.Readdir(0)
	if err != nil {
		return fmt.Errorf("failed to read migrations directory: %w", err)
	}

	for _, file := range files {
		// skipping directories and apply only .sql files
		if file.IsDir() || filepath.Ext(file.Name()) != ".sql" {
			continue
		}

		migrationFile := fmt.Sprintf("%s/%s", migrationsDir, file.Name())
		f, err := os.Open(migrationFile)
		if err != nil {
			return fmt.Errorf("failed to open migration file %s: %w", migrationFile, err)
		}
		defer f.Close()

		migrationSQL, err := io.ReadAll(f)
		if err != nil {
			return fmt.Errorf("failed to read migration file %s: %w", migrationFile, err)
		}

		// Execute the migration SQL
		_, err = dbConn.Exec(string(migrationSQL))
		if err != nil {
			return fmt.Errorf("failed to apply migration %s: %w", file.Name(), err)
		}

		log.Printf("Applied migration: %s", file.Name())
	}

	log.Println("All migrations applied successfully")
	return nil
}
