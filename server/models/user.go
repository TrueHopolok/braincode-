package models

import (
	"database/sql"
	"errors"
	"os"

	"github.com/TrueHopolok/braincode-/server/config"
	"github.com/TrueHopolok/braincode-/server/db"
)

// Deletes user form the database
func UserDelete(username string) error {
	queryfile := "delete_user.sql"
	query, err := os.ReadFile(config.Get().DBqueriesPath + queryfile)
	if err != nil {
		return err
	}

	tx, err := db.Conn.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	res, err := tx.Exec(string(query), username)
	if err != nil {
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n != 1 {
		return errors.New("invalid amount of deleted rows")
	}

	return tx.Commit()
}

// Return stats for given username
func UserFindInfo(username string) (acceptance_rate sql.NullFloat64, solved_rate sql.NullFloat64, err error) {
	queryfile := "find_user_info.sql"
	query, err := os.ReadFile(config.Get().DBqueriesPath + queryfile)
	if err != nil {
		return
	}

	tx, err := db.Conn.Begin()
	if err != nil {
		return
	}
	defer tx.Rollback()

	row := tx.QueryRow(string(query), username, username)
	if err = row.Scan(&acceptance_rate, &solved_rate); err != nil {
		return
	}

	err = tx.Commit()
	return
}

// Return salt for given username
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

// Return if user exists given username and password
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
	var ignoredUsername string
	if err := row.Scan(&ignoredUsername); err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		} else {
			return false, err
		}
	}

	return true, tx.Commit()
}

// Create a user with given username, PSH and salt
func UserCreate(username string, password, salt []byte) error {
	queryfile := "create_user.sql"
	query, err := os.ReadFile(config.Get().DBqueriesPath + queryfile)
	if err != nil {
		return err
	}

	tx, err := db.Conn.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	res, err := tx.Exec(string(query), username, password, salt)
	if err != nil {
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n != 1 {
		return errors.New("invalid amount of inserted rows")
	}

	return tx.Commit()
}
