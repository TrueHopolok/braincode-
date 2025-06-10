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

func UserDelete(w http.ResponseWriter, r *http.Request) {
	username := session.Get(r.Context()).Name
	if err := models.UserDelete(username); err != nil {
		errResp_Fatal(w, r, err)
		return
	}
	session.Logout(w)
	redirect2main(w, r, "userDelete")
}

func StatsPage(w http.ResponseWriter, r *http.Request) {
	username := session.Get(r.Context()).Name

	switch r.Header.Get("Content-Type") {
	case "text/html", "":
		ok, isenglish := langHandler(w, r)
		if !ok {
			return
		}

		acceptance_rate, solved_rate, err := models.UserFindInfo(username)
		if err != nil {
			errResp_Fatal(w, r, err)
			return
		}

		if err = views.UserFindInfo(w, username, isenglish, acceptance_rate, solved_rate); err != nil {
			errResp_Fatal(w, r, err)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
	case "application/json":
		page, err := strconv.Atoi(r.URL.Query().Get("id"))
		if err != nil || page < 0 {
			page = 0
		}

		data, err := models.SubmissionFindAll(username, page)
		if err != nil {
			errResp_Fatal(w, r, err)
			return
		}

		if _, err = w.Write(data); err != nil {
			errResp_Fatal(w, r, err)
		}
		w.Header().Set("Content-Type", "application/json")
	case "application/brainfunk":
		ssubid := r.URL.Query().Get("id")
		subid, err := strconv.Atoi(ssubid)
		if err != nil {
			http.Error(w, fmt.Sprintf("Invalid provided submission-id=%s\nWant an integer", ssubid), http.StatusNotAcceptable)
			logger.Log.Debug("req=%p submission-id=%s is not a valid integer", r, ssubid)
			return
		}
		solution, found, err := models.SubmissionFindOne(username, subid)
		if err != nil {
			errResp_Fatal(w, r, err)
			return
		} else if !found {
			http.Error(w, fmt.Sprintf("Invalid provided task-id=%d\nSuch task does not exists", subid), http.StatusNotAcceptable)
			logger.Log.Debug("req=%p task-id=%d not found", r, subid)
			return
		}

		if _, err = w.Write([]byte(solution)); err != nil {
			errResp_Fatal(w, r, err)
		}
		w.Header().Set("Content-Type", "application/brainfunk")
	default:
		denyResp_ContentTypeNotAllowed(w, r, "text/html", "application/json", "application/brainfunk")
	}
}

func RegistrationPage(w http.ResponseWriter, r *http.Request) {
	if contenttype := r.Header.Get("Content-Type"); contenttype != "" && contenttype != "text/html" {
		denyResp_ContentTypeNotAllowed(w, r, "text/html")
		return
	}

	ok, isenglish := langHandler(w, r)
	if !ok {
		return
	}

	if err := views.UserCreate(w, isenglish); err != nil {
		errResp_Fatal(w, r, err)
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
}

func UserRegister(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid login form provided", http.StatusNotAcceptable)
		logger.Log.Debug("req=%p invalid login form", r)
		return
	}
	username := r.PostFormValue("username")
	password := r.PostFormValue("password")
	if len(username) < 3 {
		http.Error(w, "Invalid login form provided\nUsername is too short", http.StatusNotAcceptable)
		logger.Log.Debug("req=%p invalid login form", r)
		return
	} else if len(password) < 8 {
		http.Error(w, "Invalid login form provided\nPassword is too short", http.StatusNotAcceptable)
		logger.Log.Debug("req=%p invalid login form", r)
		return
	}
	_, found, err := models.UserFindSalt(username)
	if err != nil {
		errResp_Fatal(w, r, err)
		return
	} else if found {
		// StatusNotAcceptable is not acceptable in this case (lol)
		http.Error(w, "User with such username exists", http.StatusNotAcceptable)
		logger.Log.Debug("req=%p trying to create same user", r)
		return
	}
	salt := SaltGen()
	if err = models.UserCreate(username, PSH(password, salt), salt); err != nil {
		errResp_Fatal(w, r, err)
		return
	}
	session.Login(session.New(username), w)
	redirect2main(w, r, "userRegister")
}

func LoginPage(w http.ResponseWriter, r *http.Request) {
	if contenttype := r.Header.Get("Content-Type"); contenttype != "" && contenttype != "text/html" {
		denyResp_ContentTypeNotAllowed(w, r, "text/html")
		return
	}

	ok, isenglish := langHandler(w, r)
	if !ok {
		return
	}

	if err := views.UserFindLogin(w, isenglish); err != nil {
		errResp_Fatal(w, r, err)
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
}

func UserLogin(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid login form provided", http.StatusNotAcceptable)
		logger.Log.Debug("req=%p invalid login form", r)
		return
	}
	username := r.PostFormValue("username")
	password := r.PostFormValue("password")
	logger.Log.Debug("req=%p username=%s password=%s", r, username, password)
	if len(username) < 3 {
		http.Error(w, "Invalid login form provided\nUsername is too short", http.StatusNotAcceptable)
		logger.Log.Debug("req=%p invalid login form", r)
		return
	} else if len(password) < 8 {
		http.Error(w, "Invalid login form provided\nPassword is too short", http.StatusNotAcceptable)
		logger.Log.Debug("req=%p invalid login form", r)
		return
	}
	salt, found, err := models.UserFindSalt(username)
	if err != nil {
		errResp_Fatal(w, r, err)
		return
	} else if !found {
		http.Error(w, "Incorrect username or password", http.StatusNotAcceptable)
		logger.Log.Debug("req=%p incorrect username or password", r)
		return
	}
	found, err = models.UserFindLogin(username, PSH(password, salt))
	if err != nil {
		errResp_Fatal(w, r, err)
		return
	} else if !found {
		http.Error(w, "Incorrect username or password", http.StatusNotAcceptable)
		logger.Log.Debug("req=%p incorrect username or password", r)
		return
	}

	ses := session.New(username)
	session.Login(ses, w)
	redirect2main(w, r, "userLogin")
}

func UserLogout(w http.ResponseWriter, r *http.Request) {
	session.Logout(w)
	redirect2main(w, r, "userLogin")
}
