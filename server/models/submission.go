package models

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/TrueHopolok/braincode-/judge"
	"github.com/TrueHopolok/braincode-/server/db"
	"github.com/TrueHopolok/braincode-/server/logger"
)

const SUBMISSIONS_AMOUNT_LIMIT = 20

type SubmissionSet struct {
	TotalAmount int
	TotalPages  int
	Rows        []SubmissionInfo
}

type SubmissionInfo struct {
	Id        int
	Timestamp string
	TaskId    sql.NullInt64
	Score     float64
}

// Return a solution for selected submission
func SubmissionFindOne(username string, subid int) (string, bool, error) {
	query, err := db.GetQuery("find_submission_one")
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

func SubmissionFindLatest(username string, taskid int) (string, bool, error) {
	query, err := db.GetQuery("find_submission_latest")
	if err != nil {
		return "", false, err
	}

	row := db.Conn.QueryRow(string(query), username, taskid, username, taskid)
	var res string
	if err := row.Scan(&res); err != nil {
		if err == sql.ErrNoRows {
			return "", false, nil
		} else {
			return "", false, err
		}

	}

	return res, true, nil
}

// Get limited amount of submissions as encoded json slice
func SubmissionFindAll(username string) ([]byte, error) {
	query, err := db.GetQuery("find_submission_all")
	if err != nil {
		return nil, err
	}

	tx, err := db.Conn.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	rows, err := tx.Query(string(query), username, SUBMISSIONS_AMOUNT_LIMIT)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rawdata SubmissionSet
	for rows.Next() {
		si := SubmissionInfo{}
		var t time.Time
		err = rows.Scan(
			&si.Id, &t,
			&si.TaskId, &si.Score,
			&rawdata.TotalAmount)
		if err != nil {
			return nil, err
		}
		si.Timestamp = t.In(time.UTC).Format("2006-01-02 15:04:05 MST")
		rawdata.Rows = append(rawdata.Rows, si)
	}
	jsondata, err := json.Marshal(rawdata)
	if err != nil {
		return nil, err
	}

	return jsondata, tx.Commit()
}

var globalJudge = judge.NewJudge(4)

// Test and get a score for a given solution and given code, then saves it into database
// Return false if solution is invalid and cannot be tested
func SubmissionCreate(username string, taskid int, solution string) (found, isvalid bool, err error) {
	findTask, err := db.GetQuery("find_task_judge")
	if err != nil {
		return false, false, err
	}

	createSubmission, err := db.GetQuery("create_submission")
	if err != nil {
		return false, false, err
	}

	updateStatus, err := db.GetQuery("update_status")
	if err != nil {
		return false, false, err
	}

	row := db.Conn.QueryRow(string(findTask), taskid)
	var rawprb []byte
	if err = row.Scan(&rawprb); err != nil {
		if err == sql.ErrNoRows {
			return false, false, nil
		} else {
			return false, false, err
		}
	}

	var prb judge.Problem
	if err = prb.UnmarshalBinary(rawprb); err != nil {
		logger.Log.Warn("task-id=%d corrupt entry", taskid)
		return true, false, err
	}
	rawverdict := globalJudge.Judge(prb, solution)
	var (
		verdict judge.Status = 0
		comment string       = ""
		score   float64      = 0
	)
	for i := range rawverdict {
		for j := range rawverdict[i] {
			verdict = rawverdict[i][j].Status
			if verdict != judge.StatusAccept {
				comment = rawverdict[i][j].Comment
				break
			}
		}
		if comment != "" {
			break
		}
	}
	if comment == "" {
		score = judge.CalculateScore(rawverdict)
	}

	tx, err := db.Conn.Begin()
	if err != nil {
		return false, false, err
	}
	defer tx.Rollback()

	res, err := tx.Exec(string(createSubmission),
		username, taskid,
		verdict, comment,
		solution, score,
		time.Now())
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

	res, err = tx.Exec(string(updateStatus), username, taskid, score, username, taskid)
	if err != nil {
		return true, true, err
	}
	n, err = res.RowsAffected()
	if err != nil {
		return true, true, err
	}
	if n < 0 || n > 1 {
		return true, true, fmt.Errorf("updated %d rows, want 0 or 1", n)
	}

	return true, true, tx.Commit()
}
