package views

import (
	"bufio"
	"net/http"

	"github.com/TrueHopolok/braincode-/server/models"
	"github.com/TrueHopolok/braincode-/server/prepared"
)

// Show 1 task page. Expects all information to be valid.
func TaskFindOne(w http.ResponseWriter, username string, isauth, isenglish bool, task models.Task) error {
	buf := bufio.NewWriter(w)
	err := prepared.Templates.ExecuteTemplate(buf, "taskpage.html", struct {
		prepared.T
	}{
		prepared.TFromBools(isenglish, isauth),
	})

	// TODO(vadim): add struct info into the page
	// anpir: left this here because still need info probably
	if err != nil {
		return err
	}
	return buf.Flush()
}

// Show problemset page. Expects all information to be valid.
func TaskFindAll(w http.ResponseWriter, username string, isauth, isenglish bool) error {
	buf := bufio.NewWriter(w)
	err := prepared.Templates.ExecuteTemplate(buf, "index.html", struct {
		Username string
		prepared.T
		Auth bool
	}{ // TODO(anpir) handle auth param in template
		Username: username,
		T:        prepared.TFromBools(isenglish, isauth),
		Auth:     isauth,
	})
	if err != nil {
		return err
	}
	return buf.Flush()
}

// Show the upload task page. Expects all information to be valid.
func TaskCreate(w http.ResponseWriter, username string, isenglish bool) error {
	buf := bufio.NewWriter(w)
	err := prepared.Templates.ExecuteTemplate(buf, "problemupload.html", struct {
		Username string
		prepared.T
	}{
		Username: username,
		T:        prepared.TFromBools(isenglish, true),
	})
	if err != nil {
		return err
	}
	return buf.Flush()
}
