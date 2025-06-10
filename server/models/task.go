package models

import (
	"cmp"
	"database/sql"
	"encoding/json"
	"errors"
	"io"

	"github.com/TrueHopolok/braincode-/judge"
	"github.com/TrueHopolok/braincode-/judge/ml"
	"github.com/TrueHopolok/braincode-/server/db"
)

type Task struct {
	General TaskInfo
	Doc     ml.Document
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
	query, err := db.GetQuery("delete_task")
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
	query, err := db.GetQuery("find_task_one.sql")
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
	var rawInfo []byte
	if err := row.Scan(
		&res.General.Id, &res.General.OwnerName,
		&res.General.TitleEn, &res.General.TitleRu,
		&rawInfo, &res.General.Score); err != nil {
		if err == sql.ErrNoRows {
			return Task{}, false, nil
		} else {
			return Task{}, false, err
		}
	}
	if err = res.Doc.UnmarshalBinary(rawInfo); err != nil {
		return Task{}, true, err
	}

	return res, true, tx.Commit()
}

// Get all task names, id and owner_id as well as amount of tasks in json
func TaskFindAll(username, search string, filter, isauth bool, page int) ([]byte, error) {
	query, err := db.GetQuery("find_task_all")
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
	if err != nil {
		return nil, err
	}

	return jsondata, tx.Commit()
}

func TaskCreate(ioDoc io.ReadCloser, username string) error {
	doc, err := ml.Parse(ioDoc)
	if err != nil {
		return err
	}
	if doc.Localizations == nil {
		return errors.New("No valid task titles were provided - Empty map")
	}
	localeEN, existsEN := doc.Localizations["en"]
	localeRU, existsRU := doc.Localizations["ru"]
	localeDEFAULT, existsDEFAULT := doc.Localizations[""]
	if !existsEN && !existsRU && !existsDEFAULT {
		return errors.New("No valid task titles were provided - No entries")
	} else if localeEN == nil && localeRU == nil && localeDEFAULT == nil {
		return errors.New("No valid task titles were provided - Nil entries")
	}

	// FIXME(anpir)
	// Previous logic is broken and this fails if no locale is provided.
	// (nil pointer dereference)
	// This is a crotch
	localeDEFAULT = cmp.Or(localeDEFAULT, new(ml.Localizable))
	localeEN = cmp.Or(localeEN, new(ml.Localizable))
	localeRU = cmp.Or(localeRU, new(ml.Localizable))

	titleDEFAULT := cmp.Or(localeDEFAULT.Name, localeEN.Name, localeRU.Name)
	if titleDEFAULT == "" {
		return errors.New("No valid task titles were provided - Zero entries")
	}

	var titleEN, titleRU string
	titleRU = cmp.Or(localeRU.Name, titleDEFAULT)
	titleEN = cmp.Or(localeEN.Name, titleDEFAULT)

	prb, err := judge.NewProblem(doc)
	if err != nil {
		return err
	}

	rawDoc, err := doc.MarshalBinary()
	if err != nil {
		return err
	}

	rawPrb, err := prb.MarshalBinary()
	if err != nil {
		return err
	}

	query, err := db.GetQuery("create_task")
	if err != nil {
		return err
	}

	tx, err := db.Conn.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	res, err := tx.Exec(string(query), username, titleEN, titleRU, rawDoc, rawPrb)
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
