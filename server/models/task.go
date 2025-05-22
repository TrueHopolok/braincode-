package models

import (
	"encoding/json"
	"os"

	"github.com/TrueHopolok/braincode-/server/config"
	"github.com/TrueHopolok/braincode-/server/db"
)

type TaskInfo struct {
	Id        int    `json:"Id"`
	Title     string `json:"Title"`
	OwnerName string `json:"OwnerName"`
}

type TasksResponse struct {
	TotalAmount int        `json:"TotalAmount"`
	Rows        []TaskInfo `json:"Rows"`
}

const TASKS_AMOUNT_LIMIT = 20

// Get all task names, id and owner_id as well as amount of tasks in json
func TaskFindAll(page int) ([]byte, error) {
	query, err := os.ReadFile(config.Get().DBqueriesPath + "task_view_all.sql")
	if err != nil {
		return nil, err
	}

	tx, err := db.Conn.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	rows, err := tx.Query(string(query), TASKS_AMOUNT_LIMIT, page*TASKS_AMOUNT_LIMIT)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rawdata TasksResponse
	for i := 0; rows.Next(); i++ {
		rawdata.Rows = append(rawdata.Rows, TaskInfo{0, "", ""})
		err = rows.Scan(&rawdata.Rows[i].Id, &rawdata.Rows[i].Title, &rawdata.Rows[i].OwnerName, &rawdata.TotalAmount)
		if err != nil {
			return nil, err
		}
	}
	jsondata, err := json.Marshal(rawdata)

	return jsondata, tx.Commit()
}
