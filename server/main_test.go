package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/TrueHopolok/braincode-/server/config"
	"github.com/TrueHopolok/braincode-/server/db"
	"github.com/TrueHopolok/braincode-/server/logger"
	"github.com/TrueHopolok/braincode-/server/prepared"
)

func TestPing(t *testing.T) {
	// TODO(anpir): fill with proper values for test
	config.OverrideConfig(t, config.Config{
		Verbose:       true,
		EnableConsole: false,
		LogFilepath:   "server.log",
		TemplatesPath: "../frontend/",
		DBuser:        "root",
		DBpass:        "root",
		DBname:        "braincode",
		DBqueriesPath: "db/queries/",
	})

	wd, _ := os.Getwd()
	fmt.Println("Working dir:", wd)

	//* Logger init
	logger.Testing()
	logger.Log.Info("[TESTING STARTED]")
	defer logger.Log.Info("[TESTING FINISHED]")

	//* Database init
	if err := db.Init(); err != nil {
		t.Error(err)
		logger.Log.Fatal("Database: connection failed; err=%s", err)
	}
	defer db.Conn.Close()

	//* Database migrate
	if err := db.Migrate(); err != nil {
		t.Error(err)
		logger.Log.Fatal("Migration: execution failed; error=%s", err)
	}

	//* Templates init
	if err := prepared.Init(); err != nil {
		t.Error(err)
		logger.Log.Fatal("Templates: initilization failed; error=%s", err)
	}

	//* HTTP init
	ts := httptest.NewServer(MuxHTTP())
	defer ts.Close()

	if _, err := http.Get(ts.URL); err != nil {
		t.Error(err)
		logger.Log.Fatal("req=GET failed; error=%s", err)
	}
}
