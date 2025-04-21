// Implements basic realization of database connection and queries.
// As well as prepared statements to use for the project.
package db

//go:generate go tool github.com/princjef/gomarkdoc/cmd/gomarkdoc -o documentation.md

import (
	"database/sql"
	"fmt"

	"github.com/TrueHopolok/braincode-/back-end/config"
	_ "github.com/ncruces/go-sqlite3/driver"
	_ "github.com/ncruces/go-sqlite3/embed"
)

// Contains pointer to sql.DB but gurantees safety of usage outside the package
type DB struct {
	*sql.DB
}

// The connection to database that is a *sql.DB type variable with limit on access it directly to avoid overwrite to nil
var Conn DB

// Open database and checks if database is reachable
func Init() error {
	var err error
	sqldb, err := sql.Open("sqlite3", fmt.Sprintf("file:%s", config.DB_FILE_PATH))
	if err != nil {
		return err
	}
	Conn = DB{sqldb}
	return Conn.Ping()
}
