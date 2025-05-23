package models

import (
	"encoding/json"
	"os"

	"github.com/TrueHopolok/braincode-/server/config"
	"github.com/TrueHopolok/braincode-/server/db"
)

type Task struct {
	General TaskInfo
	Problem []byte
}

type TaskInfo struct {
	Id        int    `json:"Id"`
	Title     string `json:"Title"`
	OwnerName string `json:"OwnerName"`
	IsSolved  bool   `json:"IsSolved"`
}

type Problemset struct {
	TotalAmount int        `json:"TotalAmount"`
	Rows        []TaskInfo `json:"Rows"`
}

const TASKS_AMOUNT_LIMIT = 20

// Get info about single task by given id and returns it as a struct
func TaskFindOne(username string, taskid int) (Task, error) {
	queryfile := "find_task_one.sql"
	query, err := os.ReadFile(config.Get().DBqueriesPath + queryfile)
	if err != nil {
		return Task{}, err
	}

	tx, err := db.Conn.Begin()
	if err != nil {
		return Task{}, err
	}
	defer tx.Rollback()

	row := tx.QueryRow(string(query), username, taskid)
	var res Task
	if err := row.Scan(); err != nil {
		return Task{}, err
	}

	return res, tx.Commit()
}

// Get all task names, id and owner_id as well as amount of tasks in json
func TaskFindAll(username string, page int) ([]byte, error) {
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

	rows, err := tx.Query(string(query), username, TASKS_AMOUNT_LIMIT, page*TASKS_AMOUNT_LIMIT)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rawdata Problemset
	for i := 0; rows.Next(); i++ {
		rawdata.Rows = append(rawdata.Rows, TaskInfo{0, "", "", false})
		err = rows.Scan(
			&rawdata.Rows[i].Id, &rawdata.Rows[i].Title,
			&rawdata.Rows[i].OwnerName, &rawdata.Rows[i].IsSolved,
			&rawdata.TotalAmount)
		if err != nil {
			return nil, err
		}
	}
	jsondata, err := json.Marshal(rawdata)

	return jsondata, tx.Commit()
}
