package views

import (
	"bufio"
	"net/http"

	"github.com/TrueHopolok/braincode-/server/models"
	"github.com/TrueHopolok/braincode-/server/prepared"
)

// Show 1 task page. Expects all information to be valid.
//
// ! TODO(vadim): wait for markleft finish
func TaskFindOne(w http.ResponseWriter, username string, isauth, isenglish bool, task models.Task) error {
	var templ string
	if isauth { // ? TODO(misha): lang and auth dependency
		if isenglish {
			templ = "TODO.html"
		} else {
			templ = "TODO.html"
		}
	} else {
		if isenglish {
			templ = "TODO.html"
		} else {
			templ = "TODO.html"
		}
	}
	templ = "taskpage.html" // TODO(vadim): delete it when other implementations will be available

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
	if isauth { // ? TODO(misha): lang and auth dependency
		if isenglish {
			templ = "index_auth.html"
		} else {
			templ = "TODO.html"
		}
	} else {
		if isenglish {
			templ = "index.html"
		} else {
			templ = "TODO.html"
		}
	}
	buf := bufio.NewWriter(w)
	err := prepared.Templates.ExecuteTemplate(buf, templ, nil) // TODO(vadim): add struct info into the page
	if err != nil {
		return err
	}
	return buf.Flush()
}

// Show the upload task page. Expects all information to be valid.
//
// ! TODO(anpir): finish markleft
func TaskCreate(w http.ResponseWriter, username string, isenglish bool) error {
	var templ string
	if isenglish { // ? TODO(misha): lang dependency
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
