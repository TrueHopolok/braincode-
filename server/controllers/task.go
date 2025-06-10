package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/TrueHopolok/braincode-/server/logger"
	"github.com/TrueHopolok/braincode-/server/models"
	"github.com/TrueHopolok/braincode-/server/session"
	"github.com/TrueHopolok/braincode-/server/views"
)

func TaskDelete(w http.ResponseWriter, r *http.Request) {
	staskid := r.URL.Query().Get("id")
	taskid, err := strconv.Atoi(staskid)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid provided task-id=%s\nWant an integer", staskid), http.StatusInternalServerError)
		logger.Log.Debug("req=%p task-id=%s is not a valid integer", r, staskid)
		return
	}

	if err := models.TaskDelete(session.Get(r.Context()).Name, taskid); err != nil {
		errResp_Fatal(w, r, err)
		return
	}
	redirect2main(w, r, "taskDelete")
}

func ProblemsPage(w http.ResponseWriter, r *http.Request) {
	ses := session.Get(r.Context())
	username := ses.Name
	isauth := !ses.IsZero()

	switch r.Header.Get("Content-Type") {
	case "text/html", "":
		ok, isenglish := langHandler(w, r)
		if !ok {
			return
		}

		if err := views.TaskFindAll(w, username, isauth, isenglish); err != nil {
			errResp_Fatal(w, r, err)
			return
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

func TaskPage(w http.ResponseWriter, r *http.Request) {
	if contenttype := r.Header.Get("Content-Type"); contenttype != "" && contenttype != "text/html" {
		denyResp_ContentTypeNotAllowed(w, r, "text/html")
		return
	}

	ok, isenglish := langHandler(w, r)
	if !ok {
		return
	}

	ses := session.Get(r.Context())
	username := ses.Name
	isauth := !ses.IsZero()

	staskid := r.URL.Query().Get("id")
	taskid, err := strconv.Atoi(staskid)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid provided task-id=%s\nWant an integer", staskid), http.StatusNotAcceptable)
		logger.Log.Debug("req=%p task-id=%s is not a valid integer", r, staskid)
		return
	}
	task, found, err := models.TaskFindOne(username, taskid)
	if err != nil {
		errResp_Fatal(w, r, err)
		return
	} else if !found {
		http.Error(w, fmt.Sprintf("Invalid provided task-id=%d\nSuch task does not exists", taskid), http.StatusNotAcceptable)
		logger.Log.Debug("req=%p task-id=%d not found", r, taskid)
		return
	}
	if err = views.TaskFindOne(w, username, isauth, isenglish, task); err != nil {
		errResp_Fatal(w, r, err)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
}

func TaskSolve(w http.ResponseWriter, r *http.Request) {
	staskid := r.URL.Query().Get("id")
	taskid, err := strconv.Atoi(staskid)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid provided task-id=%s\nWant an integer", staskid), http.StatusNotAcceptable)
		logger.Log.Debug("req=%p task-id=%s is not a valid integer", r, staskid)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid login form provided", http.StatusNotAcceptable)
		logger.Log.Debug("req=%p invalid login form", r)
		return
	}
	solution := r.PostFormValue("solution")

	found, isvalid, err := models.SubmissionCreate(session.Get(r.Context()).Name, taskid, solution)
	if err != nil {
		errResp_Fatal(w, r, err)
		return
	} else if !found {
		http.Error(w, fmt.Sprintf("Invalid provided task-id=%d\nSuch task does not exists", taskid), http.StatusNotAcceptable)
		logger.Log.Debug("req=%p task-id=%d not found", r, taskid)
		return
	} else if !isvalid {
		http.Error(w, "Given solution is invalid brainfunk code", http.StatusNotAcceptable)
		logger.Log.Debug("req=%p invalid brainfunk code", r)
		return
	}

	redirect2stats(w, r, "submitSolution")
}

func UploadPage(w http.ResponseWriter, r *http.Request) {
	if contenttype := r.Header.Get("Content-Type"); contenttype != "" && contenttype != "text/html" {
		denyResp_ContentTypeNotAllowed(w, r, "text/html")
		return
	}

	ok, isenglish := langHandler(w, r)
	if !ok {
		return
	}

	if err := views.TaskCreate(w, session.Get(r.Context()).Name, isenglish); err != nil {
		errResp_Fatal(w, r, err)
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
}

func TaskCreate(w http.ResponseWriter, r *http.Request) {
	username := session.Get(r.Context()).Name
	if err := models.TaskCreate(r.Body, username); err != nil {
		http.Error(w, fmt.Sprintf("Something went wrong while uploading the task:\n%s", err), http.StatusBadRequest)
		logger.Log.Debug("req=%p upload-err=%s ", r, err)
		return
	}
	redirect2main(w, r, "uploadTask")
}
