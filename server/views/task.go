package views

import (
	"bufio"
	"net/http"

	"github.com/TrueHopolok/braincode-/server/models"
	"github.com/TrueHopolok/braincode-/server/prepared"
)

// Show 1 task page. Expects all information to be valid.
func TaskFindOne(w http.ResponseWriter, r *http.Request, username string, isauth, isenglish bool, task models.Task) error {
	// TODO(vadim): wait for markleft finish
	templ := "taskpage.html" // lang depended
	buf := bufio.NewWriter(w)
	err := prepared.Templates.ExecuteTemplate(buf, templ, nil)
	if err != nil {
		return err
	}
	return buf.Flush()
}

// Show problemset page. Expects all information to be valid.
func TaskFindAll(w http.ResponseWriter, username string, isauth, isenglish bool) error {
	templ := "index.html" //? lang depended
	buf := bufio.NewWriter(w)
	err := prepared.Templates.ExecuteTemplate(buf, templ, struct {
		IsAuth    bool
		IsEnglish bool
		Username  string
	}{
		IsAuth:    isauth,
		IsEnglish: isenglish,
		Username:  username,
	})
	if err != nil {
		return err
	}
	return buf.Flush()
}

// Show the upload task page. Expects all information to be valid.
//
// ! TODO(anpir): finish markleft
func TaskCreate(w http.ResponseWriter, r *http.Request, username string, isenglish bool) error {
	templ := "TODO.html" //? lang depended
	buf := bufio.NewWriter(w)
	err := prepared.Templates.ExecuteTemplate(buf, templ, nil) //TODO(vadim): add struct info into the page
	if err != nil {
		return err
	}
	return buf.Flush()
}
