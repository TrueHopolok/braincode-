// Implements basic realization of database connection and queries.
// As well as prepared statements to use for the project.
package db_test

import (
	"testing"

	"github.com/TrueHopolok/braincode-/server/db"
)

func TestGetQuery(t *testing.T) {
	data, err := db.GetQuery("delete_user")
	if err != nil {
		t.Fatalf("failed to get query: %v", err)
	}
	if len(data) == 0 {
		t.Fatal("received an empty query")
	}
}
