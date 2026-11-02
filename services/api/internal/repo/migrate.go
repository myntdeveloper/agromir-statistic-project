package repo

import "database/sql"

func Migrate(db *sql.DB) error {
	return migrateBoilers(db)
}
