package models

import (
	"database/sql"
	"os"

	"github.com/TrueHopolok/braincode-/server/config"
	"github.com/TrueHopolok/braincode-/server/db"
)

func TaskFindAll(limit, page int) (*sql.Rows, error) {
	query, err := os.ReadFile(config.Get().DBFilepath + "problemset.sql")
	if err != nil {
		return nil, err
	}

	tx, err := db.Conn.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	rows, err := tx.Query(string(query), limit, page*limit)
	if err != nil {
		return nil, err
	}

	return rows, tx.Commit()
}
