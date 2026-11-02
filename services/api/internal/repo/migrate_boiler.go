package repo

import (
	"database/sql"
	_ "embed"
	"fmt"
	"log"
)

//go:embed migrations/001_boilers_schema.sql
var boilersSchemaSQL string

//go:embed migrations/002_boilers_seed.sql
var boilersSeedSQL string

func migrateBoilers(db *sql.DB) error {
	if _, err := db.Exec(boilersSchemaSQL); err != nil {
		return fmt.Errorf("boilers schema: %w", err)
	}
	return seedBoilers(db)
}

func seedBoilers(db *sql.DB) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.Exec(`LOCK TABLE warehouses IN EXCLUSIVE MODE`); err != nil {
		return fmt.Errorf("boilers seed lock: %w", err)
	}

	var seeded bool
	if err := tx.QueryRow(`SELECT EXISTS (SELECT 1 FROM warehouses)`).Scan(&seeded); err != nil {
		return fmt.Errorf("boilers seed check: %w", err)
	}
	if seeded {
		return nil
	}

	if _, err := tx.Exec(boilersSeedSQL); err != nil {
		return fmt.Errorf("boilers seed: %w", err)
	}
	if err := tx.Commit(); err != nil {
		return err
	}
	log.Println("Boilers demo data seeded")
	return nil
}
