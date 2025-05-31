package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/TrueHopolok/braincode-/server/logger"
	"github.com/TrueHopolok/braincode-/server/models"
	"github.com/TrueHopolok/braincode-/server/views"
)

func taskDelete(w http.ResponseWriter, r *http.Request) {
	username, isauth := sessionHandler(w, r)
	if !isauth {
		denyResp_NotAuthorized(w, r)
		return
	}

	staskid := r.URL.Query().Get("id")
	taskid, err := strconv.Atoi(staskid)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid provided task-id=%s\nWant an integer", staskid), 406)
		logger.Log.Debug("res=%p task-id=%s is not a valid integer", r, staskid)
		return
	}

	if err := models.TaskDelete(username, taskid); err != nil {
		errResp_Fatal(w, r, err)
		return
	}
	redirect2main(w, r, "taskDelete")
}

func getProblemset(w http.ResponseWriter, r *http.Request) {
	username, isauth := sessionHandler(w, r)

	switch r.Header.Get("Content-Type") {
	case "text/html", "":
		ok, isenglish := langHandler(w, r)
		if !ok {
			return
		} else if !isenglish {
			errResp_NotImplemented(w, r, "translation")
			return
		}

		if err := views.TaskFindAll(w, username, isauth, isenglish); err != nil {
			errResp_Fatal(w, r, err)
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
	case "application/json":
		page, err := strconv.Atoi(r.URL.Query().Get("id"))
		if err != nil || page < 0 {
			page = 0
		}

		search := r.URL.Query().Get("id")
		filter := r.URL.Query().Get("id") == "user-only"
		data, err := models.TaskFindAll(username, search, filter, isauth, page)
		if err != nil {
			errResp_Fatal(w, r, err)
			return
		}

		if _, err = w.Write(data); err != nil {
			errResp_Fatal(w, r, err)
		}
		w.Header().Set("Content-Type", "application/json")
	default:
		denyResp_ContentTypeNotAllowed(w, r, "text/html", "application/json")
	}
}

func ProblemsetPage(w http.ResponseWriter, r *http.Request) {
	logger.Log.Debug("req=%p arrived", r)
	defer logger.Log.Debug("req=%p served", r)

	switch r.Method {
	case "GET":
		getProblemset(w, r)
	case "DELETE":
		taskDelete(w, r)
	default:
		denyResp_MethodNotAllowed(w, r, "GET", "DELETE")
	}
}

func getTask(w http.ResponseWriter, r *http.Request) {
	if contenttype := r.Header.Get("Content-Type"); contenttype != "" && contenttype != "text/html" {
		denyResp_ContentTypeNotAllowed(w, r, "text/html")
		return
	}

	ok, isenglish := langHandler(w, r)
	if !ok {
		return
	} else if !isenglish {
		errResp_NotImplemented(w, r, "translation")
		return
	}

	username, isauth := sessionHandler(w, r)

	staskid := r.URL.Query().Get("id")
	taskid, err := strconv.Atoi(staskid)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid provided task-id=%s\nWant an integer", staskid), 406)
		logger.Log.Debug("res=%p task-id=%s is not a valid integer", r, staskid)
		return
	}
	task, found, err := models.TaskFindOne(username, taskid)
	if err != nil {
		errResp_Fatal(w, r, err)
		return
	} else if !found {
		http.Error(w, fmt.Sprintf("Invalid provided task-id=%d\nSuch task does not exists", taskid), 406)
		logger.Log.Debug("res=%p task-id=%d not found", r, taskid)
		return
	}
	// TODO(VADIM): add markleft handler to info and then pass it to view
	if err = views.TaskFindOne(w, r, username, isauth, isenglish, task); err != nil {
		errResp_Fatal(w, r, err)
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
}

func submitSolution(w http.ResponseWriter, r *http.Request) {
	username, isauth := sessionHandler(w, r)
	if !isauth {
		denyResp_NotAuthorized(w, r)
		return
	}

	staskid := r.URL.Query().Get("id")
	taskid, err := strconv.Atoi(staskid)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid provided task-id=%s\nWant an integer", staskid), 406)
		logger.Log.Debug("res=%p task-id=%s is not a valid integer", r, staskid)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid login form provided", 406)
		logger.Log.Debug("res=%p invalid login form", r)
		return
	}
	solution := r.PostFormValue("solution")

	found, isvalid, err := models.SubmissionCreate(username, taskid, solution)
	if err != nil {
		errResp_Fatal(w, r, err)
		return
	} else if !found {
		http.Error(w, fmt.Sprintf("Invalid provided task-id=%d\nSuch task does not exists", taskid), 406)
		logger.Log.Debug("res=%p task-id=%d not found", r, taskid)
		return
	} else if !isvalid {
		http.Error(w, fmt.Sprintf("Given solution is invalid brainfunk code"), 406)
		logger.Log.Debug("req=%p invalid brainfunk code", r)
		return
	}

	redirect2stats(w, r, "submitSolution")
}

func TaskPage(w http.ResponseWriter, r *http.Request) {
	logger.Log.Debug("req=%p arrived", r)
	defer logger.Log.Debug("req=%p served", r)

	switch r.Method {
	case "GET":
		getTask(w, r)
	case "POST":
		submitSolution(w, r)
	default:
		denyResp_MethodNotAllowed(w, r, "GET", "POST")
	}
}

func getUpload(w http.ResponseWriter, r *http.Request, username string) {
	if contenttype := r.Header.Get("Content-Type"); contenttype != "" && contenttype != "text/html" {
		denyResp_ContentTypeNotAllowed(w, r, "text/html")
		return
	}

	ok, isenglish := langHandler(w, r)
	if !ok {
		return
	} else if !isenglish {
		errResp_NotImplemented(w, r, "translation")
		return
	}

	errResp_NotImplemented(w, r, "getUpload")
}

func uploadTask(w http.ResponseWriter, r *http.Request, username string) {
	errResp_NotImplemented(w, r, "uploadTask")
}

func UploadPage(w http.ResponseWriter, r *http.Request) {
	logger.Log.Debug("req=%p arrived", r)
	defer logger.Log.Debug("req=%p served", r)

	username, isauth := sessionHandler(w, r)
	if !isauth {
		denyResp_NotAuthorized(w, r)
		return
	}

	switch r.Method {
	case "GET":
		getUpload(w, r, username)
	case "POST":
		uploadTask(w, r, username)
	default:
		denyResp_MethodNotAllowed(w, r, "GET", "POST")
	}
}
