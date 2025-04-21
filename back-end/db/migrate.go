package db

import (
	"embed"
	"fmt"
	"io"
	"io/fs"
	"slices"

	"github.com/TrueHopolok/braincode-/back-end/logger"
)

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

	for _, entry := range entries {
		f, err := migrations.Open(entry)
		if err != nil {
			return fmt.Errorf("migration %s: cannot open embedded file: %w", entry, err)
		}

		if err := runMigration(f); err != nil {
			return fmt.Errorf("migration %s: %w", entry, err)
		}
	}

	return nil
}

func runMigration(f fs.File) error {
	stat, err := f.Stat()
	if err != nil {
		// Unreachable, embed.FS file State() never fails.
		panic(fmt.Errorf("migration: cannot stat file: %w", err))
	}

	logger.Log.Debug("Migration: found = %s", stat.Name())
	data, err := io.ReadAll(f)
	if err != nil {
		// Unreachable, embed.FS file Read() never fails.
		panic(fmt.Errorf("migration %s: cannot read file: %w", stat.Name(), err))
	}

	if _, err := Conn.Exec(string(data)); err != nil {
		return fmt.Errorf("migration %s: %w", stat.Name(), err)
	}

	return nil
}
