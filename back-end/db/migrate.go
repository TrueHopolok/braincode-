package db

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"os"
)

var MIGRATIONS_DIR_PARH = flag.String("db_migrations", "back-end/db/migrations/", "Provide a directory with all migartion queries")

func Migrate(migration_name ...string) error {
	if conn == nil {
		return errors.New("database is not initialized")
	}
	for _, mname := range migration_name {
		mname = fmt.Sprintf("%s%s%s", *MIGRATIONS_DIR_PARH, mname, ".sql")
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
