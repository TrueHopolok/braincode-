package views

import (
	"bufio"
	"net/http"

	"github.com/TrueHopolok/braincode-/judge/ml"
	"github.com/TrueHopolok/braincode-/server/models"
	"github.com/TrueHopolok/braincode-/server/prepared"
)

// Show 1 task page. Expects all information to be valid.
func TaskFindOne(w http.ResponseWriter, username string, isauth, isenglish bool, task models.Task, previousSolution string) error {
	buf := bufio.NewWriter(w)
	t := prepared.TFromBools(isenglish, isauth)
	err := prepared.Templates.ExecuteTemplate(buf, "taskpage.html", struct {
		prepared.T
		Document ml.TemplatableDocument
		Solution string
	}{
		t,
		task.Doc.Templatable(t.Lang),
		previousSolution,
	})
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
	}{
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
