package models

import (
	"database/sql"
	"encoding/json"
	"errors"
	"os"
	"time"

	"github.com/TrueHopolok/braincode-/server/config"
	"github.com/TrueHopolok/braincode-/server/db"
)

const SUBMISSIONS_AMOUNT_LIMIT = 20

type SubmissionSet struct {
	TotalAmount int              `json:"TotalAmount"`
	Rows        []SubmissionInfo `json:"Rows"`
}

type SubmissionInfo struct {
	Id        int           `json:"Id"`
	Timestamp time.Time     `json:"Timestamp"`
	TaskId    sql.NullInt64 `json:"TaskId"`
	Score     float64       `json:"Score"`
}

// Return a solution for selected submission
func SubmissionFindOne(username string, subid int) (string, bool, error) {
	queryfile := "find_submission_one.sql"
	query, err := os.ReadFile(config.Get().DBqueriesPath + queryfile)
	if err != nil {
		return "", false, err
	}

	tx, err := db.Conn.Begin()
	if err != nil {
		return "", false, err
	}
	defer tx.Rollback()

	row := tx.QueryRow(string(query), username, subid)
	var res string
	if err := row.Scan(&res); err != nil {
		if err == sql.ErrNoRows {
			return "", false, nil
		} else {
			return "", false, err
		}

	}

	return res, true, tx.Commit()
}

// Get limited amount of submissions as encoded json slice
func SubmissionFindAll(username string, page int) ([]byte, error) {
	queryfile := "find_submission_all.sql"
	query, err := os.ReadFile(config.Get().DBqueriesPath + queryfile)
	if err != nil {
		return nil, err
	}

	tx, err := db.Conn.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	rows, err := tx.Query(string(query), username, SUBMISSIONS_AMOUNT_LIMIT, SUBMISSIONS_AMOUNT_LIMIT*page)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rawdata SubmissionSet
	for i := 0; rows.Next(); i++ {
		rawdata.Rows = append(rawdata.Rows, SubmissionInfo{})
		err = rows.Scan(
			&rawdata.Rows[i].Id, &rawdata.Rows[i].Timestamp,
			&rawdata.Rows[i].TaskId, &rawdata.Rows[i].Score,
			&rawdata.TotalAmount)
		if err != nil {
			return nil, err
		}
	}
	jsondata, err := json.Marshal(rawdata)

	return jsondata, tx.Commit()
}

// Test and get a score for a given solution and given code, then saves it into database
// Return false if solution is invalid and cannot be tested
func SubmissionCreate(username string, taskid int, solution string) (found, isvalid bool, err error) {
	queryfile := "find_task_judge.sql"
	query1, err := os.ReadFile(config.Get().DBqueriesPath + queryfile)
	if err != nil {
		return false, false, err
	}

	queryfile = "create_submission.sql"
	query2, err := os.ReadFile(config.Get().DBqueriesPath + queryfile)
	if err != nil {
		return false, false, err
	}

	queryfile = "update_status.sql"
	query3, err := os.ReadFile(config.Get().DBqueriesPath + queryfile)
	if err != nil {
		return false, false, err
	}

	tx, err := db.Conn.Begin()
	if err != nil {
		return false, false, err
	}
	defer tx.Rollback()

	row := tx.QueryRow(string(query1), taskid)
	var rawproblem []byte
	if err = row.Scan(&rawproblem); err != nil {
		if err == sql.ErrNoRows {
			return false, false, nil
		} else {
			return false, false, err
		}
	}
	if err = tx.Commit(); err != nil {
		return true, false, err
	}

	// TODO(vadim): JUDGE to recieve score and check if solution is valid
	score := 1.0

	res, err := tx.Exec(string(query2), username, taskid, solution, score, time.Now())
	if err != nil {
		return true, true, err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return true, true, err
	}
	if n != 1 {
		return true, true, errors.New("invalid amount of inserted rows")
	}

	res, err = tx.Exec(string(query3), username, taskid, score, username, taskid)
	if err != nil {
		return true, true, err
	}
	n, err = res.RowsAffected()
	if err != nil {
		return true, true, err
	}
	if n != 1 {
		return true, true, errors.New("invalid amount of updated rows")
	}

	return true, true, tx.Commit()
}
