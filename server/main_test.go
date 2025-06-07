package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/TrueHopolok/braincode-/server/config"
	"github.com/TrueHopolok/braincode-/server/db"
	"github.com/TrueHopolok/braincode-/server/prepared"
)

func TestPing(t *testing.T) {
	// TODO(anpir): fill with proper values for test
	config.OverrideConfig(t, config.Config{
		Verbose:       true,
		LogFilepath:   "",
		TemplatesPath: "",
		DBuser:        "",
		DBpass:        "",
		DBname:        "",
		DBqueriesPath: "",
	})

	//* Database init
	if err := db.Init(); err != nil {
		t.Error(err)
	}
	defer db.Conn.Close()

	//* Database migrate
	if err := db.Migrate(); err != nil {
		t.Error(err)
	}

	//* Templates init
	if err := prepared.Init(); err != nil {
		t.Error(err)
	}

	//* HTTP init
	ts := httptest.NewServer(MuxHTTP())
	defer ts.Close()

	//* Ping test itself
	if _, err := http.Get(ts.URL); err != nil {
		t.Error(err)
	}
}
