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

func userDelete(w http.ResponseWriter, r *http.Request, username string) {
	if err := models.UserDelete(username); err != nil {
		errResp_Fatal(w, r, err)
		return
	}
	w.Header().Del("Session")
	redirect2main(w, r, "userDelete")
}

func getStats(w http.ResponseWriter, r *http.Request, username string) {
	switch r.Header.Get("Content-Type") {
	case "text/html", "":
		ok, isenglish := langHandler(w, r)
		if !ok {
			return
		} else if !isenglish {
			errResp_NotImplemented(w, r, "translation")
			return
		}

		acceptance_rate, solved_rate, err := models.UserFindInfo(username)
		if err != nil {
			errResp_Fatal(w, r, err)
			return
		}

		if err = views.UserFindInfo(w, r, username, isenglish, acceptance_rate, solved_rate); err != nil {
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
			http.Error(w, fmt.Sprintf("Invalid provided submission-id=%s\nWant an integer", ssubid), 406)
			logger.Log.Debug("res=%p submission-id=%s is not a valid integer", r, ssubid)
			return
		}
		solution, found, err := models.SubmissionFindOne(username, subid)
		if err != nil {
			errResp_Fatal(w, r, err)
			return
		} else if !found {
			http.Error(w, fmt.Sprintf("Invalid provided task-id=%d\nSuch task does not exists", subid), 406)
			logger.Log.Debug("res=%p task-id=%d not found", r, subid)
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

func StatsPage(w http.ResponseWriter, r *http.Request) {
	logger.Log.Debug("req=%p arrived", r)
	defer logger.Log.Debug("req=%p served", r)

	username, isauth := sessionHandler(w, r)
	if !isauth {
		denyResp_NotAuthorized(w, r)
		return
	}

	switch r.Method {
	case "GET":
		getStats(w, r, username)
	case "DELETE":
		userDelete(w, r, username)
	default:
		denyResp_MethodNotAllowed(w, r, "GET", "DELETE")
	}
}

func getRegistration(w http.ResponseWriter, r *http.Request) {
	if contenttype := r.Header.Get("Content-Type"); contenttype != "" && contenttype != "text/html" {
		denyResp_ContentTypeNotAllowed(w, r, "text/html")
		return
	}

	_, isauth := sessionHandler(w, r)
	if isauth {
		denyResp_DenyAuthorized(w, r)
		return
	}

	ok, isenglish := langHandler(w, r)
	if !ok {
		return
	} else if !isenglish {
		errResp_NotImplemented(w, r, "translation")
		return
	}

	if err := views.UserCreate(w, r, isenglish); err != nil {
		errResp_Fatal(w, r, err)
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
}

func userRegister(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid login form provided", 406)
		logger.Log.Debug("res=%p invalid login form", r)
		return
	}
	username := r.PostFormValue("username")
	password := r.PostFormValue("password")
	if len(username) < 3 {
		http.Error(w, "Invalid login form provided\nUsername is too short", 406)
		logger.Log.Debug("res=%p invalid login form", r)
		return
	} else if len(password) < 8 {
		http.Error(w, "Invalid login form provided\nPassword is too short", 406)
		logger.Log.Debug("res=%p invalid login form", r)
		return
	}
	_, found, err := models.UserFindSalt(username)
	if err != nil {
		errResp_Fatal(w, r, err)
		return
	} else if !found {
		http.Error(w, "User with such username exists", 406)
		logger.Log.Debug("res=%p trying to create same user", r)
		return
	}
	salt := SaltGen()
	if err = models.UserCreate(username, PSH(password, salt), salt); err != nil {
		errResp_Fatal(w, r, err)
		return
	}
	w.Header().Set("Session", session.New(username).CreateJWT())
	redirect2main(w, r, "userRegister")
}

func RegistrationPage(w http.ResponseWriter, r *http.Request) {
	logger.Log.Debug("req=%p arrived", r)
	defer logger.Log.Debug("req=%p served", r)

	switch r.Method {
	case "GET":
		getRegistration(w, r)
	case "POST":
		userRegister(w, r)
	default:
		denyResp_MethodNotAllowed(w, r, "GET", "POST")
	}
}

func getLogin(w http.ResponseWriter, r *http.Request) {
	if contenttype := r.Header.Get("Content-Type"); contenttype != "" && contenttype != "text/html" {
		denyResp_ContentTypeNotAllowed(w, r, "text/html")
		return
	}

	_, isauth := sessionHandler(w, r)
	if isauth {
		denyResp_DenyAuthorized(w, r)
		return
	}

	ok, isenglish := langHandler(w, r)
	if !ok {
		return
	} else if !isenglish {
		errResp_NotImplemented(w, r, "translation")
		return
	}

	if err := views.UserFindLogin(w, r, isenglish); err != nil {
		errResp_Fatal(w, r, err)
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
}

func userAuth(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid login form provided", 406)
		logger.Log.Debug("res=%p invalid login form", r)
		return
	}
	username := r.PostFormValue("username")
	password := r.PostFormValue("password")
	if len(username) < 3 {
		http.Error(w, "Invalid login form provided\nUsername is too short", 406)
		logger.Log.Debug("res=%p invalid login form", r)
		return
	} else if len(password) < 8 {
		http.Error(w, "Invalid login form provided\nPassword is too short", 406)
		logger.Log.Debug("res=%p invalid login form", r)
		return
	}
	salt, found, err := models.UserFindSalt(username)
	if err != nil {
		errResp_Fatal(w, r, err)
		return
	} else if !found {
		http.Error(w, "Incorrect username or password", 406)
		logger.Log.Debug("res=%p incorrect username or password", r)
		return
	}
	found, err = models.UserFindLogin(username, PSH(password, salt))
	if err != nil {
		errResp_Fatal(w, r, err)
		return
	} else if !found {
		http.Error(w, "Incorrect username or password", 406)
		logger.Log.Debug("res=%p incorrect username or password", r)
		return
	}
	w.Header().Set("Session", session.New(username).CreateJWT())
	redirect2main(w, r, "userLogin")
}

func LoginPage(w http.ResponseWriter, r *http.Request) {
	logger.Log.Debug("req=%p arrived", r)
	defer logger.Log.Debug("req=%p served", r)

	switch r.Method {
	case "GET":
		getLogin(w, r)
	case "POST":
		userAuth(w, r)
	default:
		denyResp_MethodNotAllowed(w, r, "GET", "POST")
	}
}
