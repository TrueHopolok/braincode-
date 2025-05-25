package controllers

import (
	"net/http"

	"github.com/TrueHopolok/braincode-/server/logger"
	"github.com/TrueHopolok/braincode-/server/models"
	"github.com/TrueHopolok/braincode-/server/session"
	"github.com/TrueHopolok/braincode-/server/views"
)

func StatsPage(w http.ResponseWriter, r *http.Request) {
	logger.Log.Debug("req=%p arrived", r)
	defer logger.Log.Debug("req=%p served", r)

	if r.Method != "GET" {
		errResponseMethodNotAllowed(w, r, "GET")
		return
	}

	username, isauth := sessionHandler(w, r)
	if !isauth {
		errResponseNotAuthorized(w, r)
		return
	}

	if contenttype := r.Header.Get("Content-Type"); contenttype != "" && contenttype != "text/html" {
		errResponseContentTypeNotAllowed(w, r, "text/html")
		return
	}

	ok, isenglish := langHandler(w, r)
	if !ok {
		return
	} else if !isenglish {
		errResponseNotImplemented(w, r, "translation")
		return
	}

	userinfo, err := models.UserFindInfo(username)
	if err != nil {
		errResponseFatal(w, r, err)
		return
	}

	if err = views.UserFindInfo(userinfo); err != nil {
		errResponseFatal(w, r, err)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
}

func getpageRegister(w http.ResponseWriter, r *http.Request) {
	if contenttype := r.Header.Get("Content-Type"); contenttype != "" && contenttype != "text/html" {
		errResponseContentTypeNotAllowed(w, r, "text/html")
		return
	}

	ok, isenglish := langHandler(w, r)
	if !ok {
		return
	} else if !isenglish {
		errResponseNotImplemented(w, r, "translation")
		return
	}

	username, isauth := sessionHandler(w, r)

	if err := views.UserCreate(w, r, username, isauth, isenglish); err != nil {
		errResponseFatal(w, r, err)
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
		errResponseFatal(w, r, err)
		return
	} else if !found {
		http.Error(w, "User with such username exists", 406)
		logger.Log.Debug("res=%p trying to create same user", r)
		return
	}
	salt := SaltGen()
	if err = models.UserCreate(username, PSH(password, salt), salt); err != nil {
		errResponseFatal(w, r, err)
		return
	}
	w.Header().Set("Session", session.New(username).CreateJWT())
	w.WriteHeader(204)
}

func RegistrationPage(w http.ResponseWriter, r *http.Request) {
	logger.Log.Debug("req=%p arrived", r)
	defer logger.Log.Debug("req=%p served", r)

	switch r.Method {
	case "GET":
		getpageRegister(w, r)
	case "POST":
		userRegister(w, r)
	default:
		errResponseMethodNotAllowed(w, r, "GET", "POST")
	}
}

func getpageLogin(w http.ResponseWriter, r *http.Request) {
	if contenttype := r.Header.Get("Content-Type"); contenttype != "" && contenttype != "text/html" {
		errResponseContentTypeNotAllowed(w, r, "text/html")
		return
	}

	ok, isenglish := langHandler(w, r)
	if !ok {
		return
	} else if !isenglish {
		errResponseNotImplemented(w, r, "translation")
		return
	}

	username, isauth := sessionHandler(w, r)

	if err := views.UserFindLogin(w, r, username, isauth, isenglish); err != nil {
		errResponseFatal(w, r, err)
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
		errResponseFatal(w, r, err)
		return
	} else if !found {
		http.Error(w, "Incorrect username or password", 406)
		logger.Log.Debug("res=%p incorrect username or password", r)
		return
	}
	found, err = models.UserFindLogin(username, PSH(password, salt))
	if err != nil {
		errResponseFatal(w, r, err)
		return
	} else if !found {
		http.Error(w, "Incorrect username or password", 406)
		logger.Log.Debug("res=%p incorrect username or password", r)
		return
	}
	w.Header().Set("Session", session.New(username).CreateJWT())
	w.WriteHeader(204)
}

func LoginPage(w http.ResponseWriter, r *http.Request) {
	logger.Log.Debug("req=%p arrived", r)
	defer logger.Log.Debug("req=%p served", r)

	switch r.Method {
	case "GET":
		getpageLogin(w, r)
	case "POST":
		userAuth(w, r)
	default:
		errResponseMethodNotAllowed(w, r, "GET", "POST")
	}
}
