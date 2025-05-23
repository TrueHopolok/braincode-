package models

import (
	"database/sql"
	"os"

	"github.com/TrueHopolok/braincode-/server/config"
	"github.com/TrueHopolok/braincode-/server/db"
)

func UserFindSalt(username string) ([]byte, bool, error) {
	queryfile := "find_user_salt.sql"
	query, err := os.ReadFile(config.Get().DBqueriesPath + queryfile)
	if err != nil {
		return nil, false, err
	}

	tx, err := db.Conn.Begin()
	if err != nil {
		return nil, false, err
	}
	defer tx.Rollback()

	row := tx.QueryRow(string(query), username)
	var salt []byte
	if err := row.Scan(&salt); err != nil {
		if err == sql.ErrNoRows {
			return nil, false, nil
		} else {
			return nil, false, err
		}
	}

	return salt, true, tx.Commit()
}

func UserFindLogin(username string, password []byte) (bool, error) {
	queryfile := "find_user_login.sql"
	query, err := os.ReadFile(config.Get().DBqueriesPath + queryfile)
	if err != nil {
		return false, err
	}

	tx, err := db.Conn.Begin()
	if err != nil {
		return false, err
	}
	defer tx.Rollback()

	row := tx.QueryRow(string(query), username, password)
	if err := row.Scan(); err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		} else {
			return false, err
		}
	}

	return true, tx.Commit()
}
