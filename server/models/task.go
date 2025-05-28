package models

import (
	"database/sql"
	"encoding/json"
	"errors"
	"os"

	"github.com/TrueHopolok/braincode-/server/config"
	"github.com/TrueHopolok/braincode-/server/db"
)

type Task struct {
	General TaskInfo
	Info    []byte
}

type TaskInfo struct {
	Id        int             `json:"Id"`
	TitleEn   string          `json:"TitleEn"`
	TitleRu   string          `json:"TitleRu"`
	OwnerName string          `json:"OwnerName"`
	Score     sql.NullFloat64 `json:"Score"`
}

type Problemset struct {
	TotalAmount int        `json:"TotalAmount"`
	Rows        []TaskInfo `json:"Rows"`
}

const TASKS_AMOUNT_LIMIT = 20

// Deletes task from the database
func TaskDelete(username string, taskid int) error {
	queryfile := "delete_task.sql"
	query, err := os.ReadFile(config.Get().DBqueriesPath + queryfile)
	if err != nil {
		return err
	}

	tx, err := db.Conn.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	res, err := tx.Exec(string(query), username, taskid)
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

// Get info about single task by given id and returns it as a struct
func TaskFindOne(username string, taskid int) (Task, bool, error) {
	queryfile := "find_task_one.sql"
	query, err := os.ReadFile(config.Get().DBqueriesPath + queryfile)
	if err != nil {
		return Task{}, false, err
	}

	tx, err := db.Conn.Begin()
	if err != nil {
		return Task{}, false, err
	}
	defer tx.Rollback()

	row := tx.QueryRow(string(query), username, taskid)
	var res Task
	if err := row.Scan(
		&res.General.Id, &res.General.OwnerName,
		&res.General.TitleEn, &res.General.TitleRu,
		&res.Info, &res.General.Score); err != nil {
		if err == sql.ErrNoRows {
			return Task{}, false, nil
		} else {
			return Task{}, false, err
		}
	}

	return res, true, tx.Commit()
}

// Get all task names, id and owner_id as well as amount of tasks in json
func TaskFindAll(username, search string, filter, isauth bool, page int) ([]byte, error) {
	queryfile := "find_task_all.sql"
	query, err := os.ReadFile(config.Get().DBqueriesPath + queryfile)
	if err != nil {
		return nil, err
	}

	tx, err := db.Conn.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	rows, err := tx.Query(string(query),
		username,
		search,
		search,
		username,
		!(filter && isauth),
		TASKS_AMOUNT_LIMIT, TASKS_AMOUNT_LIMIT*page)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rawdata Problemset
	for i := 0; rows.Next(); i++ {
		rawdata.Rows = append(rawdata.Rows, TaskInfo{})
		err = rows.Scan(
			&rawdata.Rows[i].Id, &rawdata.Rows[i].TitleEn, &rawdata.Rows[i].TitleRu,
			&rawdata.Rows[i].OwnerName, &rawdata.Rows[i].Score,
			&rawdata.TotalAmount)
		if err != nil {
			return nil, err
		}
	}
	jsondata, err := json.Marshal(rawdata)

	return jsondata, tx.Commit()
}
