package db

import (
	"database/sql"
	"fmt"
	"testing"
	"unicode"

	"github.com/TrueHopolok/braincode-/server/config"
)

func InitTesting(t *testing.T) {
	t.Helper()
	sqldb, err := sql.Open("mysql",
		fmt.Sprintf("%s:%s@/",
			config.Get().DBuser,
			config.Get().DBpass,
		))
	if err != nil {
		t.Fatal(t)
	}

	for _, r := range config.Get().DBname {
		if !((unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_') && r <= 127) {
			t.Errorf("db name %q may not be a valid sql identifier, which may brake tests", config.Get().DBname)
			break
		}
	}

	q := []string{
		`DROP DATABASE IF EXISTS %s;`,
		`CREATE DATABASE %s;`,
	}

	for _, qq := range q {
		if _, err := sqldb.Exec(fmt.Sprintf(qq, config.Get().DBname)); err != nil {
			t.Fatalf("failed to wipe previous database: err = %v", err)
		}
	}

	if err := sqldb.Close(); err != nil {
		t.Fatalf("could not close sql con: err = %v", err)
	}

	if err := Init(); err != nil {
		t.Fatalf("db initialization failed: err = %v", err)
	}
}
