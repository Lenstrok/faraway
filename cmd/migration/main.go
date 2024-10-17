package main

import (
	"fmt"
	"log"

	"vio_coding_challenge/internal/config"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

const dbDriverName = "postgres"

func main() {
	cfg := config.MigrationConfig{}

	err := config.NewConfig(&cfg)
	if err != nil {
		log.Fatalf("Failed to get config: %v", err)
	}

	// For production run we should add lifetime & conn pool size.
	db, err := sqlx.Connect(dbDriverName, cfg.DB.GetDSN())
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}

	if err = migrateUp(db); err != nil {
		log.Fatalf("Failed to migrate up: %v", err)
	}

	log.Printf("Migration completed successfully.")
}

func migrateUp(db *sqlx.DB) error {
	driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("failed to get driver: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance("file://././migration", dbDriverName, driver)
	if err != nil {
		return fmt.Errorf("failed to get migration instance: %w", err)
	}

	log.Printf("DB Mirgation log: %v", m.Up())

	return nil
}
