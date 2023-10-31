package data

import (
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"os"
)

func MigrateUp() {
	dbURL := os.Getenv("dbURL")
	migrationPath := os.Getenv("migrationPath")

	// Create a new instance of migrate.
	m, err := migrate.New(migrationPath, dbURL)
	if err != nil {
		fmt.Println("Error creating migrate instance:", err)
		os.Exit(1)
	}

	// Migrations to the last version.
	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		fmt.Println("Error applying migrations:", err)
		os.Exit(1)
	}

	fmt.Println("Migrations applied successfully.")
}
