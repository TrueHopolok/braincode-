package db

import (
	"database/sql"
	"errors"
	"flag"
	"fmt"

	_ "github.com/ncruces/go-sqlite3/driver"
	_ "github.com/ncruces/go-sqlite3/embed"
)

var DB_FILE_PATH = flag.String("db_location", "back-end/db/db.db", "Provide location for database file location")

var conn *sql.DB

func Init() error {
	var err error
	conn, err = sql.Open("sqlite3", fmt.Sprintf("file:%s", *DB_FILE_PATH))
	if err != nil {
		return err
	}
	return conn.Ping()
}

func Execute(query_path string, data ...any) error {
	if conn == nil {
		return errors.New("database is not initialized")
	}
	// TODO: add functionality for executing queries
	return nil
}
