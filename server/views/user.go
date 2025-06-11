package views

import (
	"bufio"
	"database/sql"
	"net/http"

	"github.com/TrueHopolok/braincode-/server/prepared"
)

// Show user's stats via template with prepared section to handle fetch request of submitions. Expects all information to be valid.
func UserFindInfo(w http.ResponseWriter, username string, isenglish bool, acceptanceRate, solvedRate sql.NullFloat64) error {
	if !acceptanceRate.Valid {
		acceptanceRate.Float64 = 0
	}
	if !solvedRate.Valid {
		solvedRate.Float64 = 0
	}

	ar := 0.0
	if acceptanceRate.Valid {
		ar = acceptanceRate.Float64 * 100
	}

	sr := 0.0
	if solvedRate.Valid {
		sr = solvedRate.Float64 * 100
	}

	buf := bufio.NewWriter(w)
	err := prepared.Templates.ExecuteTemplate(buf, "userpage.html", struct {
		Username       string
		AcceptanceRate float64
		SolvedRate     float64
		prepared.T
	}{
		Username:       username,
		AcceptanceRate: ar,
		SolvedRate:     sr,
		T:              prepared.TFromBools(isenglish, true),
	})
	if err != nil {
		return err
	}
	return buf.Flush()
}

// Show login page. Expects all information to be valid.
func UserFindLogin(w http.ResponseWriter, isenglish bool) error {
	buf := bufio.NewWriter(w)
	err := prepared.Templates.ExecuteTemplate(buf, "login.html", struct {
		prepared.T
	}{
		T: prepared.TFromBools(isenglish, false),
	})
	if err != nil {
		return err
	}
	return buf.Flush()
}

// Show registration page. Expects all information to be valid.
func UserCreate(w http.ResponseWriter, isenglish bool) error {
	buf := bufio.NewWriter(w)
	err := prepared.Templates.ExecuteTemplate(buf, "registration.html", struct {
		prepared.T
	}{
		prepared.TFromBools(isenglish, false),
	})
	if err != nil {
		return err
	}
	return buf.Flush()
}
