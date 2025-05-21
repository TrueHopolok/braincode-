package views

import (
	"bufio"
	"database/sql"
	"net/http"

	"github.com/TrueHopolok/braincode-/server/prepared"
)

func Problemset(w http.ResponseWriter, templ, username string, isauth bool, rows *sql.Rows) error {
	// TODO(vadim): to the finished view add all information
	buf := bufio.NewWriter(w)
	err := prepared.Templates.ExecuteTemplate(buf, templ, nil)
	if err != nil {
		return err
	}
	return buf.Flush()
}
