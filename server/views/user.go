package views

import (
	"bufio"
	"database/sql"
	"net/http"

	"github.com/TrueHopolok/braincode-/server/prepared"
)

// Show user's stats via template with prepared section to handle fetch request of submitions. Expects all information to be valid.
func UserFindInfo(w http.ResponseWriter, username string, isenglish bool, acceptance_rate, solved_rate sql.NullFloat64) error {
	var templ string // ? TODO(misha): lang dependency
	if isenglish {
		templ = "TODO.html"
	} else {
		templ = "TODO.html"
	}

	buf := bufio.NewWriter(w)
	err := prepared.Templates.ExecuteTemplate(buf, templ, nil) // TODO(vadim): add struct info into the page
	if err != nil {
		return err
	}
	return buf.Flush()
}

// Show login page. Expects all information to be valid.
func UserFindLogin(w http.ResponseWriter, isenglish bool) error {
	var templ string // ? TODO(misha): lang dependency
	if isenglish {
		templ = "login.html"
	} else {
		templ = "TODO.html"
	}

	buf := bufio.NewWriter(w)
	err := prepared.Templates.ExecuteTemplate(buf, templ, nil) // No struct info is needed
	if err != nil {
		return err
	}
	return buf.Flush()
}

// Show registration page. Expects all information to be valid.
func UserCreate(w http.ResponseWriter, isenglish bool) error {
	var templ string // ? TODO(misha): lang dependency
	if isenglish {
		templ = "registration.html"
	} else {
		templ = "TODO.html"
	}

	buf := bufio.NewWriter(w)
	err := prepared.Templates.ExecuteTemplate(buf, templ, nil) // No struct info is needed
	if err != nil {
		return err
	}
	return buf.Flush()
}
