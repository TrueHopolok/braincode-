package views

import (
	"bufio"
	"net/http"

	"github.com/TrueHopolok/braincode-/server/models"
	"github.com/TrueHopolok/braincode-/server/prepared"
)

// Show 1 task page. Expects all information to be valid.
func TaskFindOne(w http.ResponseWriter, username string, isauth, isenglish bool, task models.Task) error {
	var templ string
	if isenglish {
		templ = "taskpage.html"
	} else {
		templ = "taskpage_ru.html"
	}

	buf := bufio.NewWriter(w)
	err := prepared.Templates.ExecuteTemplate(buf, templ, nil) // TODO(vadim): add struct info into the page
	if err != nil {
		return err
	}
	return buf.Flush()
}

// Show problemset page. Expects all information to be valid.
func TaskFindAll(w http.ResponseWriter, username string, isauth, isenglish bool) error {
	var templ string
	if isauth {
		if isenglish {
			templ = "index_auth.html"
		} else {
			templ = "index_auth_ru.html"
		}
	} else {
		if isenglish {
			templ = "index.html"
		} else {
			templ = "index_ru.html"
		}
	}
	buf := bufio.NewWriter(w)
	err := prepared.Templates.ExecuteTemplate(buf, templ, struct {
		Username string
	}{
		Username: username,
	})
	if err != nil {
		return err
	}
	return buf.Flush()
}

// Show the upload task page. Expects all information to be valid.
func TaskCreate(w http.ResponseWriter, username string, isenglish bool) error {
	var templ string
	if isenglish {
		templ = "problemupload.html"
	} else {
		templ = "problemupload_ru.html"
	}

	buf := bufio.NewWriter(w)
	err := prepared.Templates.ExecuteTemplate(buf, templ, struct {
		Username string
	}{
		Username: username,
	})
	if err != nil {
		return err
	}
	return buf.Flush()
}
