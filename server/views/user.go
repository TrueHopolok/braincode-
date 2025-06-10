package views

import (
	"bufio"
	"database/sql"
	"math"
	"net/http"

	"github.com/TrueHopolok/braincode-/server/prepared"
)

// Show user's stats via template with prepared section to handle fetch request of submitions. Expects all information to be valid.
func UserFindInfo(w http.ResponseWriter, username string, isenglish bool, acceptance_rate, solved_rate sql.NullFloat64) error {
	var templ string
	if isenglish {
		templ = "userpage.html"
	} else {
		templ = "userpage_ru.html"
	}

	if !acceptance_rate.Valid {
		acceptance_rate.Float64 = 0
	}
	if !solved_rate.Valid {
		solved_rate.Float64 = 0
	}

	buf := bufio.NewWriter(w)
	err := prepared.Templates.ExecuteTemplate(buf, templ, struct {
		Username       string
		AcceptanceRate float64
		SolvedRate     float64
	}{
		Username:       username,
		AcceptanceRate: math.Round(acceptance_rate.Float64*1000) / 10,
		SolvedRate:     math.Round(solved_rate.Float64*1000) / 10,
	})
	if err != nil {
		return err
	}
	return buf.Flush()
}

// Show login page. Expects all information to be valid.
func UserFindLogin(w http.ResponseWriter, isenglish bool) error {
	var templ string
	if isenglish {
		templ = "login.html"
	} else {
		templ = "login_ru.html"
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
	var templ string
	if isenglish {
		templ = "registration.html"
	} else {
		templ = "registration_ru.html"
	}

	buf := bufio.NewWriter(w)
	err := prepared.Templates.ExecuteTemplate(buf, templ, nil) // No struct info is needed
	if err != nil {
		return err
	}
	return buf.Flush()
}
