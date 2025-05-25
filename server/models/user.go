package models

import (
	"database/sql"
	"errors"
	"os"

	"github.com/TrueHopolok/braincode-/server/config"
	"github.com/TrueHopolok/braincode-/server/db"
)

type User struct {
}

func UserFindInfo(username string) (User, error) {
	queryfile := "find_user_info.sql"
	query, err := os.ReadFile(config.Get().DBqueriesPath + queryfile)
	if err != nil {
		return User{}, err
	}

	tx, err := db.Conn.Begin()
	if err != nil {
		return User{}, err
	}
	defer tx.Rollback()

	row := tx.QueryRow(string(query), username)
	var res User
	if err := row.Scan(); err != nil {
		return User{}, err
	}

	return res, tx.Commit()
}

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
