package models

import (
	"database/sql"
	"encoding/json"
	"errors"
	"time"

	"github.com/TrueHopolok/braincode-/judge"
	"github.com/TrueHopolok/braincode-/server/db"
	"github.com/TrueHopolok/braincode-/server/logger"
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

// Get limited amount of submissions as encoded json slice
func SubmissionFindAll(username string, page int) ([]byte, error) {
	query, err := db.GetQuery("find_submission_all")
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
	if err != nil {
		return nil, err
	}

	return jsondata, tx.Commit()
}

// Test and get a score for a given solution and given code, then saves it into database
// Return false if solution is invalid and cannot be tested
func SubmissionCreate(username string, taskid int, solution string) (found, isvalid bool, err error) {
	query1, err := db.GetQuery("find_task_judge")
	if err != nil {
		return false, false, err
	}

	query2, err := db.GetQuery("create_submission")
	if err != nil {
		return false, false, err
	}

	query3, err := db.GetQuery("update_status")
	if err != nil {
		return false, false, err
	}

	tx, err := db.Conn.Begin()
	if err != nil {
		return false, false, err
	}
	defer tx.Rollback()

	row := tx.QueryRow(string(query1), taskid)
	var rawprb []byte
	if err = row.Scan(&rawprb); err != nil {
		if err == sql.ErrNoRows {
			return false, false, nil
		} else {
			return false, false, err
		}
	}
	if err = tx.Commit(); err != nil {
		return true, false, err
	}

	var prb judge.Problem
	if err = prb.UnmarshalBinary(rawprb); err != nil {
		logger.Log.Warn("task-id=%d corrupt entry", taskid)
		return true, false, err
	}
	jdg := judge.NewJudge(4)
	rawverdict := jdg.Judge(prb, solution)
	if err = jdg.Close(); err != nil {
		logger.Log.Fatal("judge error=%s", err)
	}
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

	res, err := tx.Exec(string(query2),
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
