package goose

import (
	"database/sql"
	"fmt"
)

// Assuming a dialect-specific provider interface exists for locking
// In a real scenario, this would interact with the database driver's locking mechanism.
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
	// 1. Acquire Lock (Dialect-specific implementation required)
	// For the purpose of this fix, we assume a Lock() function exists.
	unlock, err := Lock(db)
	if err != nil {
		return fmt.Errorf("failed to acquire migration lock: %w", err)
	}
	defer unlock()

	// 2. Get Current Version AFTER acquiring lock
	currentVersion, err := GetDBVersion(db)
	if err != nil {
		return fmt.Errorf("failed to get current database version: %w", err)
	}

	// 3. Filter and Execute
	for _, m := range migrations {
		if m.Version > currentVersion {
			if err := runMigration(db, m, true); err != nil {
				return err
			}
		}
	}
	return nil
}