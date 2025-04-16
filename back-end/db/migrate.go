package db

import (
	"fmt"
	"os"
)

const MIGRATIONS_DIR_PARH = "back-end/migrations/"

func Migrate(migration_path string) error {
	migration_path = fmt.Sprintf("%s%s%s", MIGRATIONS_DIR_PARH, migration_path, "sql")
	if _, err := os.Stat(migration_path); err != nil {
		return err
	}
	// TODO: sql execution logic
	return nil
}