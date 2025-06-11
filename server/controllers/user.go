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

	if err = views.UserFindInfo(w, username, isenglish, acceptance_rate, solved_rate); err != nil {
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
		page, err := strconv.Atoi(r.URL.Query().Get("page"))
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
		http.Error(w, "Invalid login form provided", http.StatusBadRequest)
		logger.Log.Debug("req=%p invalid registration form", r)
		return
	}
	username := r.PostFormValue("username")
	password := r.PostFormValue("password")
	if len(username) < 3 {
		http.Error(w, "Invalid login form provided\nUsername is too short", http.StatusBadRequest)
		logger.Log.Debug("req=%p invalid login form", r)
		return
	} else if len(password) < 8 {
		http.Error(w, "Invalid login form provided\nPassword is too short", http.StatusBadRequest)
		logger.Log.Debug("req=%p invalid login form", r)
		return
	}
	_, found, err := models.UserFindSalt(username)
	if err != nil {
		errResp_Fatal(w, r, err)
		return
	} else if found {
		http.Error(w, "User with such username exists", http.StatusBadRequest)
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

func UserChangePassword(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid password change form provided", http.StatusBadRequest)
	}

	username := session.Get(r.Context()).Name
	passOld := r.FormValue("current_password")
	passNew := r.FormValue("new_password")
	passConfirm := r.FormValue("confirm_password")

	// double check user auth
	if ok, err := authValid(username, passOld); err != nil {
		errResp_Fatal(w, r, err)
		return
	} else if !ok {
		http.Error(w, "Old password does not match", http.StatusBadRequest)
		logger.Log.Debug("req=%p incorrect old password", r)
		return
	}

	// validate data

	if passNew == passOld {
		http.Error(w, "New password must differ from old one", http.StatusBadRequest)
		logger.Log.Debug("req=%p duplicate new password", r)
		return
	}

	if passConfirm != passNew {
		http.Error(w, "Passwords do not match", http.StatusBadRequest)
		logger.Log.Debug("req=%p passwords do not match", r)
		return
	}

	// ok!
	salt := SaltGen()
	if err := models.UserChangePassword(username, PSH(passNew, salt), salt); err != nil {
		errResp_Fatal(w, r, err)
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

	if err := views.UserFindLogin(w, isenglish); err != nil {
		errResp_Fatal(w, r, err)
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
}

func UserLogin(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid login form provided", http.StatusBadRequest)
		logger.Log.Debug("req=%p invalid login form", r)
		return
	}
	username := r.PostFormValue("username")
	password := r.PostFormValue("password")
	logger.Log.Debug("req=%p username=%s password=%s", r, username, password)
	if len(username) < 3 {
		http.Error(w, "Invalid login form provided\nUsername is too short", http.StatusBadRequest)
		logger.Log.Debug("req=%p invalid login form", r)
		return
	} else if len(password) < 8 {
		http.Error(w, "Invalid login form provided\nPassword is too short", http.StatusBadRequest)
		logger.Log.Debug("req=%p invalid login form", r)
		return
	}

	if ok, err := authValid(username, password); err != nil {
		errResp_Fatal(w, r, err)
		return
	} else if !ok {
		http.Error(w, "Incorrect username or password", http.StatusBadRequest)
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
