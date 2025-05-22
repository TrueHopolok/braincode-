package views

import (
	"bufio"
	"net/http"

	"github.com/TrueHopolok/braincode-/server/prepared"
)

func TaskViewAll(w http.ResponseWriter, templ, username string, isauth bool) error {
	// TODO(vadim): to the finished view add all information
	buf := bufio.NewWriter(w)
	err := prepared.Templates.ExecuteTemplate(buf, templ, nil)
	if err != nil {
		return err
	}
	return buf.Flush()
}
