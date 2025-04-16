package db

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"os"
)

// Contain migration files directory
var MIGRATIONS_DIR_PATH = flag.String("db_migrations", "back-end/db/migrations/", "Provide a directory with all migartion queries")

// Execute all given queries if they exists as files.
// Only requirement is to provide name of the file, not the path nor extenstion.
func Migrate(migration_name ...string) error {
	if conn == nil {
		return errors.New("database is not initialized")
	}
	for _, mname := range migration_name {
		mname = fmt.Sprintf("%s%s%s", *MIGRATIONS_DIR_PATH, mname, ".sql")
		mfile, err := os.Open(mname)
		if err != nil {
			return err
		}
		var buf bytes.Buffer
		_, err = buf.ReadFrom(mfile)
		mfile.Close()
		if err != nil {
			return err
		}
		_, err = conn.Exec(buf.String())
		if err != nil {
			return err
		}
	}
	return nil
}
