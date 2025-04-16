package db

import (
	"database/sql"
	"fmt"

	_ "github.com/ncruces/go-sqlite3/driver"
	_ "github.com/ncruces/go-sqlite3/embed"
)

const DB_FILE_PATH = "back-end/db/db.db"

func Version() {
	var version string
	db, _ := sql.Open("sqlite3", fmt.Sprintf("file:%s", DB_FILE_PATH))
	db.QueryRow(`SELECT sqlite_version()`).Scan(&version)
	fmt.Println(version)
}

/*

TODO:
Database connection, sql query execution.
In other files all queries are checked and queries prepared for execution.
Meaning this file only for logic without that much error handling

*/
