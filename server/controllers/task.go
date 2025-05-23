package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/TrueHopolok/braincode-/server/logger"
	"github.com/TrueHopolok/braincode-/server/models"
	"github.com/TrueHopolok/braincode-/server/views"
)

func Problemset(w http.ResponseWriter, r *http.Request) {
	logger.Log.Debug("req=%p arrived", r)
	defer logger.Log.Debug("req=%p served", r)

	if r.Method != "GET" {
		errResponseMethodNotAllowed(w, r, "GET")
		return
	}

	ses, isauth := sessionHandler(w, r)

	switch r.Header.Get("Content-Type") {
	case "text/html", "":
		ok, isenglish := langHandler(w, r)
		if !ok {
			return
		} else if !isenglish {
			errResponseNotImplemented(w, r, "translation")
			return
		}

		if err := views.TaskFindAll(w, ses.Name, isauth, isenglish); err != nil {
			errResponseFatal(w, r, err)
		}
	case "application/json":
		page, err := strconv.Atoi(r.Header.Get("Page"))
		if err != nil || page < 0 {
			page = 0
		}

		search := r.Header.Get("Search")
		filter := r.Header.Get("Filter") == "user-only"
		data, err := models.TaskFindAll(ses.Name, search, filter, isauth, page)
		if err != nil {
			errResponseFatal(w, r, err)
			return
		}

		if _, err = w.Write(data); err != nil {
			errResponseFatal(w, r, err)
		}
	default:
		errResponseContentTypeNotAllowed(w, r, "text/html", "application/json")
	}
}

func getTaskPage(w http.ResponseWriter, r *http.Request) {
	ses, isauth := sessionHandler(w, r)

	ok, isenglish := langHandler(w, r)
	if !ok {
		return
	} else if !isenglish {
		errResponseNotImplemented(w, r, "translation")
		return
	}

	if r.Header.Get("Content-Type") != "text/html" && r.Header.Get("Content-Type") != "" {
		errResponseContentTypeNotAllowed(w, r, "text/html")
		return
	}

	staskid := r.URL.Query().Get("id")
	taskid, err := strconv.Atoi(staskid)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid provided task-id=%s\nWant an integer", staskid), 406)
		logger.Log.Debug("res=%p task-id=%s is not a valid integer", r, staskid)
		return
	}
	task, found, err := models.TaskFindOne(ses.Name, taskid)
	if err != nil {
		errResponseFatal(w, r, err)
		return
	} else if !found {
		http.Error(w, fmt.Sprintf("Invalid provided task-id=%d\nSuch task does not exists", taskid), 406)
		logger.Log.Debug("res=%p task-id=%d not found", r, taskid)
		return
	}
	if err = views.TaskFindOne(w, r, ses.Name, isauth, isenglish, task); err != nil {
		errResponseFatal(w, r, err)
	}
}

func postTaskPage(w http.ResponseWriter, r *http.Request) {
	errResponseNotImplemented(w, r, "postTaskPage")
	// TODO(vadim): add post req handler
	// Check if auth
	// Get a submission
	// Save into to the database
}

func TaskPage(w http.ResponseWriter, r *http.Request) {
	logger.Log.Debug("req=%p arrived", r)
	defer logger.Log.Debug("req=%p served", r)

	switch r.Method {
	case "GET":
		getTaskPage(w, r)
	case "POST":
		postTaskPage(w, r)
	default:
		errResponseMethodNotAllowed(w, r, "GET", "POST")
	}
}
