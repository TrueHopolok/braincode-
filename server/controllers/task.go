package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/TrueHopolok/braincode-/server/logger"
	"github.com/TrueHopolok/braincode-/server/models"
	"github.com/TrueHopolok/braincode-/server/views"
)

const TASKS_ON_1_PAGE = 20

type TaskInfo struct {
	Id    int    `json:"id"`
	Title string `json:"title"`
}

type TasksResponse struct {
	TotalAmount int        `json:"totalAmount"`
	Rows        []TaskInfo `json:"rows"`
}

func Problemset(w http.ResponseWriter, r *http.Request) {
	logger.Log.Debug("req=%p arrived", r)
	defer logger.Log.Debug("req=%p served", r)

	if r.Method != "GET" {
		http.Error(w, fmt.Sprintf("Method=%s is not allowed\nAllowed=GET", r.Method), 405)
		w.Header().Set("Allow", "GET")
		logger.Log.Debug("req=%p method=%s is not allowed", r, r.Method)
		return
	}

	ses, isauth := sessionHandler(w, r)
	contenttype := r.Header.Get("Content-Type")

	switch contenttype {
	case "text/html", "":
		templ := "index.html"
		lang := r.URL.Query().Get("lang")
		if lang == "ru" {
			// templ = "index_ru.html" // TODO(vadim): switch to this when it will be done
			notImplemented(w, r, "translation is not available yet")
			return
		} else if lang != "" && lang != "en" {
			http.Error(w, "Such language selection is not allowed", 406)
			logger.Log.Debug("req=%p lang=%s is not allowed", r, lang)
			return
		}

		if err := views.TaskViewAll(w, templ, ses.Name, isauth, lang); err != nil {
			http.Error(w, "Failed to write into the response body", 500)
			logger.Log.Error("req=%p failed; error=%s", r, err)
		}
	case "application/json":
		page, err := strconv.Atoi(r.Header.Get("Page"))
		if err != nil || page < 0 {
			page = 0
		}
		rows, err := models.TaskFindAll(TASKS_ON_1_PAGE, page)
		if err != nil {
			http.Error(w, "Failed to write into the response body", 500)
			logger.Log.Error("req=%p failed; error=%s", r, err)
			return
		}
		var data TasksResponse
		for i := 0; rows.Next(); i++ {
			data.Rows = append(data.Rows, TaskInfo{0, ""})
			rows.Scan(&data.Rows[i].Id, &data.Rows[i].Title, &data.TotalAmount)
		}
		logger.Log.Debug("%v", data)
		raw, err := json.Marshal(data)
		if err != nil {
			http.Error(w, "Failed to write into the response body", 500)
			logger.Log.Error("req=%p failed; error=%s", r, err)
			return
		}
		if _, err = w.Write(raw); err != nil {
			http.Error(w, "Failed to write into the response body", 500)
			logger.Log.Error("req=%p failed; error=%s", r, err)
		}
	default:
		http.Error(w, fmt.Sprintf("Content-Type=%s is not allowed\nAllowed=text/html application/json", contenttype), 406)
		logger.Log.Debug("req=%p Content-Type=%s is not allowed", r, contenttype)
	}
}
