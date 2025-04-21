package db

import (
	"embed"
	"errors"
	"fmt"
	"io/fs"
	"slices"

	"github.com/TrueHopolok/braincode-/back-end/logger"
)

const migrationVersionTable = "MigrationVersion"

// migrations contains embedded migration files.
// This filesystem should only contain sql migration files and directories.
// It may be arbitrarily nested.
// Migrations will be applied in alphabetic order.
//
//go:embed migrations/*.sql
var migrations embed.FS

// Migrate executes all embedded migrations.
func Migrate() error {
	var entries []string
	if err := fs.WalkDir(migrations, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			entries = append(entries, path)
		}
		return nil
	}); err != nil {
		return fmt.Errorf("migration: cannot read embedded directory: %w", err)
	}

	slices.Sort(entries)

	var version string
	{ // lookup database version
		var dirty bool
		if err := Conn.QueryRow("SELECT version, dirty FROM "+migrationVersionTable+";").Scan(&version, &dirty); err != nil {
			// There is no way to check whether a table exists while staying drive-agnostic.
			// Because of that, we assume that version table does not exist on any failure.
			version = ""
			dirty = false

			if _, err := Conn.Exec("CREATE TABLE " + migrationVersionTable + " (version TEXT, dirty BOOLEAN, pk INTEGER PRIMARY KEY);"); err != nil {
				return fmt.Errorf("migration: cannot create meta table: %w", err)
			}
			if _, err := Conn.Exec("INSERT INTO " + migrationVersionTable + " VALUES ('', FALSE, 0);"); err != nil {
				return fmt.Errorf("migration: cannot create meta table: %w", err)
			}
		}

		if dirty {
			return errors.New("migration: database is dirty, this can happen if something went horribly wrong mid-migration, fix the database manually")
		}
	}

	i, found := slices.BinarySearch(entries, version)
	if !found && version != "" {
		return fmt.Errorf("migration: database version %v is not known", version)
	}

	for _, entry := range entries[i:] {
		data, err := migrations.ReadFile(entry)
		if err != nil {
			return fmt.Errorf("migration %s: cannot open embedded file: %w", entry, err)
		}

		if err := runMigration(entry, data); err != nil {
			return fmt.Errorf("migration %s: %w", entry, err)
		}
	}

	return nil
}

func runMigration(version string, data []byte) error {
	logger.Log.Debug("Migration: found = %s", version)

	tx, err := Conn.Begin()
	if err != nil {
		return fmt.Errorf("migration %s: cannot start transaction: %w", version, err)
	}
	defer tx.Rollback()

	if _, err := tx.Exec("UPDATE " + migrationVersionTable + " SET dirty=TRUE;"); err != nil {
		return fmt.Errorf("migration %s: cannot write meta table: %w", version, err)
	}

	if _, err := tx.Exec(string(data)); err != nil {
		return fmt.Errorf("migration %s: %w", version, err)
	}

	if _, err := tx.Exec("UPDATE "+migrationVersionTable+" SET dirty=FALSE, version=?;", version); err != nil {
		return fmt.Errorf("migration %s: cannot write meta table: %w", version, err)
	}

	return tx.Commit()
}
