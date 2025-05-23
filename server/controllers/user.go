package controllers

import (
	"net/http"

	"github.com/TrueHopolok/braincode-/server/logger"
	"github.com/TrueHopolok/braincode-/server/models"
	"github.com/TrueHopolok/braincode-/server/session"
	"github.com/TrueHopolok/braincode-/server/views"
)

func getpageRegister(w http.ResponseWriter, r *http.Request) {
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

	if err := views.UserCreate(w, r, ses.Name, isauth, isenglish); err != nil {
		errResponseFatal(w, r, err)
	}
}

func userRegister(w http.ResponseWriter, r *http.Request) {
	errResponseNotImplemented(w, r, "userRegister")
	// TODO(vadim): add post req handler
	// Check data and type of request
	// Return either error in registration/login or user data as sesion token
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

	if err := views.UserFindLogin(w, r, ses.Name, isauth, isenglish); err != nil {
		errResponseFatal(w, r, err)
	}
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
	errResponseNotImplemented(w, r, "success")
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
