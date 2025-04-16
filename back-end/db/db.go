// Implements basic realization of database connection and queries.
// As well as prepared statements to use for the project.
package db

//go:generate go tool github.com/princjef/gomarkdoc/cmd/gomarkdoc -o documentation.md

import (
	"database/sql"
	"flag"
	"fmt"

	_ "github.com/ncruces/go-sqlite3/driver"
	_ "github.com/ncruces/go-sqlite3/embed"
)

// Contain path to database file
var DB_FILE_PATH = flag.String("db_location", "back-end/db/db.db", "Provide location for database file location")

// Contains pointer to sql.DB but gurantees safety of usage outside the package
type DB struct {
	*sql.DB
}

// The connection to database that is a *sql.DB type variable with limit on access it directly to avoid overwrite to nil
var Conn DB

// Open database and checks if database is reachable
func Init() error {
	var err error
	sqldb, err := sql.Open("sqlite3", fmt.Sprintf("file:%s", *DB_FILE_PATH))
	if err != nil {
		return err
	}
	Conn = DB{sqldb}
	return Conn.Ping()
}
