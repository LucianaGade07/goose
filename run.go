package goose

import (
	"database/sql"
	"fmt"
)

func runMigration(db *sql.DB, migration *Migration, direction bool) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	if err := migration.Run(tx); err != nil {
		tx.Rollback()
		return fmt.Errorf("migration %v failed: %w", migration.Version, err)
	}

	return tx.Commit()
}

func Up(db *sql.DB, migrations []*Migration) error {
	for _, m := range migrations {
		if err := runMigration(db, m, true); err != nil {
			return err // Halt execution immediately
		}
	}
	return nil
}