package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/TrueHopolok/braincode-/server/logger"
	"github.com/TrueHopolok/braincode-/server/models"
	"github.com/TrueHopolok/braincode-/server/prepared"
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

func ProfilePage(w http.ResponseWriter, r *http.Request) {
	username := session.Get(r.Context()).Name
	ok, isenglish := langHandler(w, r)
	if !ok {
		return
	}

	acceptance_rate, solved_rate, err := models.UserFindInfo(username)
	if err != nil {
		errResp_Fatal(w, r, err)
		return
	}

	errorcode := prepared.T{}.Request(r).ErrCode

	if err = views.UserFindInfo(w, username, isenglish, acceptance_rate, solved_rate, errorcode); err != nil {
		errResp_Fatal(w, r, err)
		return
	}
}

func SubmissionsAPI(w http.ResponseWriter, r *http.Request) {
	username := session.Get(r.Context()).Name

	if r.URL.Query().Has("id") {
		// get a singular submition
		ssubid := r.URL.Query().Get("id")
		subid, err := strconv.Atoi(ssubid)
		if err != nil {
			http.Error(w, fmt.Sprintf("Invalid provided submission-id=%s\nWant an integer", ssubid), http.StatusBadRequest)
			logger.Log.Debug("req=%p submission-id=%s is not a valid integer", r, ssubid)
			return
		}
		solution, found, err := models.SubmissionFindOne(username, subid)
		if err != nil {
			errResp_Fatal(w, r, err)
			return
		} else if !found {
			http.Error(w, fmt.Sprintf("Invalid provided task-id=%d\nSuch task does not exists", subid), http.StatusBadRequest)
			logger.Log.Debug("req=%p task-id=%d not found", r, subid)
			return
		}

		w.Header().Set("Content-Type", "application/brainfunk") // lol

		if _, err = w.Write([]byte(solution)); err != nil {
			errResp_Fatal(w, r, err)
		}

	} else {
		// get list of all submissions
		data, err := models.SubmissionFindAll(username)
		if err != nil {
			errResp_Fatal(w, r, err)
			return
		}

		if _, err = w.Write(data); err != nil {
			errResp_Fatal(w, r, err)
		}
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

	if err := views.UserCreate(w, isenglish, prepared.T{}.Request(r).ErrCode); err != nil {
		errResp_Fatal(w, r, err)
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
}

func UserRegister(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		redirectError(w, r, 1) // bad request
		logger.Log.Debug("req=%p invalid registration form", r)
		return
	}
	username := r.PostFormValue("username")
	password := r.PostFormValue("password")
	if len(username) < 3 {
		redirectError(w, r, 2) // short username
		logger.Log.Debug("req=%p invalid login form", r)
		return
	} else if len(password) < 8 {
		redirectError(w, r, 3) // short password
		logger.Log.Debug("req=%p invalid login form", r)
		return
	}
	_, found, err := models.UserFindSalt(username)
	if err != nil {
		redirectError(w, r, 4) // internal error
		return
	} else if found {
		redirectError(w, r, 5) // already exists
		logger.Log.Debug("req=%p trying to create same user", r)
		return
	}
	salt := SaltGen()
	if err = models.UserCreate(username, PSH(password, salt), salt); err != nil {
		redirectError(w, r, 4) // internal error
		return
	}
	session.Login(session.New(username), w)
	redirect2main(w, r, "userRegister")
}

func UserChangePassword(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid password change form provided", http.StatusBadRequest)
		redirectError(w, r, 1) // bad request
		return
	}

	username := session.Get(r.Context()).Name
	passOld := r.FormValue("current_password")
	passNew := r.FormValue("new_password")
	passConfirm := r.FormValue("confirm_password")

	// double check user auth
	if ok, err := authValid(username, passOld); err != nil {
		redirectError(w, r, 2) // internal error
		return
	} else if !ok {
		redirectError(w, r, 3) // bad old pass
		logger.Log.Debug("req=%p incorrect old password", r)
		return
	}

	// validate data

	if passNew == passOld {
		redirectError(w, r, 4) // old pass same as new
		logger.Log.Debug("req=%p duplicate new password", r)
		return
	}

	if passConfirm != passNew {
		redirectError(w, r, 5) // confirmation does not match
		logger.Log.Debug("req=%p passwords do not match", r)
		return
	}

	if len(passNew) < 8 {
		redirectError(w, r, 6) // new password too short
		logger.Log.Debug("req=%p password too short", r)
		return
	}

	// ok!
	salt := SaltGen()
	if err := models.UserChangePassword(username, PSH(passNew, salt), salt); err != nil {
		redirectError(w, r, 2) // internal error
		return
	}

	redirect2stats(w, r, "userChanePassword")
}

func LoginPage(w http.ResponseWriter, r *http.Request) {
	if contenttype := r.Header.Get("Content-Type"); contenttype != "" && contenttype != "text/html" {
		denyResp_ContentTypeNotAllowed(w, r, "text/html")
		logger.Log.Debug("req=%p invalid password change form form", r)
		return
	}

	ok, isenglish := langHandler(w, r)
	if !ok {
		return
	}

	errcode := prepared.T{}.Request(r).ErrCode
	if err := views.UserFindLogin(w, isenglish, errcode); err != nil {
		errResp_Fatal(w, r, err)
	}
}

func UserLogin(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		redirectError(w, r, 1) // bad form
		logger.Log.Debug("req=%p invalid login form", r)
		return
	}
	username := r.PostFormValue("username")
	password := r.PostFormValue("password")
	logger.Log.Debug("req=%p username=%s password=%s", r, username, password)
	if len(username) < 3 {
		redirectError(w, r, 2) // short login
		logger.Log.Debug("req=%p invalid login form", r)
		return
	} else if len(password) < 8 {
		redirectError(w, r, 3) // short pass
		logger.Log.Debug("req=%p invalid login form", r)
		return
	}

	if ok, err := authValid(username, password); err != nil {
		redirectError(w, r, 4) // internal error
		return
	} else if !ok {
		redirectError(w, r, 5) // wrong username/password
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

func authValid(user, pass string) (bool, error) {
	salt, found, err := models.UserFindSalt(user)
	if err != nil {
		return false, err
	}
	if !found {
		return false, nil
	}

	found, err = models.UserFindLogin(user, PSH(pass, salt))
	if err != nil {
		return false, err
	}
	if !found {
		return false, nil
	}

	return true, nil
}
