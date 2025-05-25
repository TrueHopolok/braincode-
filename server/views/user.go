package views

import (
	"bufio"
	"net/http"

	"github.com/TrueHopolok/braincode-/server/models"
	"github.com/TrueHopolok/braincode-/server/prepared"
)

func UserFindInfo(w http.ResponseWriter, r *http.Request, userinfo models.User) error {
	// TODO(vadim): to the finished view add all information
	templ := "TODO.html" // lang depended
	buf := bufio.NewWriter(w)
	err := prepared.Templates.ExecuteTemplate(buf, templ, nil)
	if err != nil {
		return err
	}
	return buf.Flush()
}

func UserFindLogin(w http.ResponseWriter, r *http.Request, username string, isauth, isenglish bool) error {
	// TODO(vadim): to the finished view add all information
	templ := "login.html" // lang depended
	buf := bufio.NewWriter(w)
	err := prepared.Templates.ExecuteTemplate(buf, templ, nil)
	if err != nil {
		return err
	}
	return buf.Flush()
}

func UserCreate(w http.ResponseWriter, r *http.Request, username string, isauth, isenglish bool) error {
	// TODO(vadim): to the finished view add all information
	templ := "registration.html" // lang depended
	buf := bufio.NewWriter(w)
	err := prepared.Templates.ExecuteTemplate(buf, templ, nil)
	if err != nil {
		return err
	}
	return buf.Flush()
}
