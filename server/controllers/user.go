package controllers

import (
	"net/http"

	"github.com/TrueHopolok/braincode-/server/logger"
	"github.com/TrueHopolok/braincode-/server/views"
)

func getRegistrationPage(w http.ResponseWriter, r *http.Request) {
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

func postRegistrationPage(w http.ResponseWriter, r *http.Request) {
	errResponseNotImplemented(w, r, "postLoginPage")
	// TODO(vadim): add post req handler
	// Check data and type of request
	// Return either error in registration/login or user data as sesion token
}

func RegistrationPage(w http.ResponseWriter, r *http.Request) {
	logger.Log.Debug("req=%p arrived", r)
	defer logger.Log.Debug("req=%p served", r)

	switch r.Method {
	case "GET":
		getRegistrationPage(w, r)
	case "POST":
		postRegistrationPage(w, r)
	default:
		errResponseMethodNotAllowed(w, r, "GET", "POST")
	}
}

func getLoginPage(w http.ResponseWriter, r *http.Request) {
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

func postLoginPage(w http.ResponseWriter, r *http.Request) {
	errResponseNotImplemented(w, r, "postLoginPage")
	// TODO(vadim): add post req handler
	// Check data and type of request
	// Return either error in registration/login or user data as sesion token
}

func LoginPage(w http.ResponseWriter, r *http.Request) {
	logger.Log.Debug("req=%p arrived", r)
	defer logger.Log.Debug("req=%p served", r)

	switch r.Method {
	case "GET":
		getLoginPage(w, r)
	case "POST":
		postLoginPage(w, r)
	default:
		errResponseMethodNotAllowed(w, r, "GET", "POST")
	}
}
