package models

import (
	"database/sql"
	"errors"

	"github.com/TrueHopolok/braincode-/server/db"
)

// Deletes user form the database
func UserDelete(username string) error {
	query, err := db.GetQuery("delete_user")
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
	query, err := db.GetQuery("find_user_info")
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
	query, err := db.GetQuery("find_user_salt")
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
func UserFindLogin(username string, psh []byte) (bool, error) {
	query, err := db.GetQuery("find_user_login")
	if err != nil {
		return false, err
	}

	tx, err := db.Conn.Begin()
	if err != nil {
		return false, err
	}
	defer tx.Rollback()

	row := tx.QueryRow(string(query), username, psh)
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
func UserCreate(username string, psh, salt []byte) error {
	query, err := db.GetQuery("create_user")
	if err != nil {
		return err
	}

	tx, err := db.Conn.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	res, err := tx.Exec(string(query), username, psh, salt)
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

func UserChangePassword(username string, psh, salt []byte) error {
	query, err := db.GetQuery("change_password")
	if err != nil {
		return err
	}

	tx, err := db.Conn.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	res, err := tx.Exec(string(query), psh, salt, username)
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
