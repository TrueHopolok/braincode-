package main

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/TrueHopolok/braincode-/server/config"
	"github.com/TrueHopolok/braincode-/server/db"
	"github.com/TrueHopolok/braincode-/server/logger"
	"github.com/TrueHopolok/braincode-/server/prepared"
)

// ! Call only inside test functions !
//
// [InitBackend] should be called at the beggining of each test.
//
// Initialize everything needed for server:
//   - Config,
//   - Logger,
//   - Database (with migrations),
//   - Templates;
//
// Requires calling closing functions manually:
//   - [db.Conn.Close()],
//   - [logger.Log.Info("[TESTING FINISHED]")];
func InitBackend(t *testing.T) {
	if !testing.Testing() {
		panic("InitBackend called outside of a test")
	}

	//* Config init for testing
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

	//* Logger init
	logger.Testing()
	logger.Log.Info("[TESTING STARTED]")

	//* Database init
	if err := db.Init(); err != nil {
		t.Error(err)
		logger.Log.Error("Database: connection failed; err=%s", err)
		logger.Log.Info("[TESTING FINISHED]")
		return
	}

	//* Database migrate
	if err := db.Migrate(); err != nil {
		t.Error(err)
		logger.Log.Error("Migration: execution failed; error=%s", err)
		db.Conn.Close()
		logger.Log.Info("[TESTING FINISHED]")
		return
	}

	//* Templates init
	if err := prepared.Init(); err != nil {
		t.Error(err)
		logger.Log.Error("Templates: initilization failed; error=%s", err)
		db.Conn.Close()
		logger.Log.Info("[TESTING FINISHED]")
		return
	}
}

// ! Call only inside test functions !
//
// [MustRequest] should be called as a wrapper of [http.NewRequest] function that returns [http.Request].
// Will fail the test and panic in case of an error.
func MustRequest(t *testing.T, method, url string, body io.Reader) *http.Request {
	if !testing.Testing() {
		panic("MustRequest called outside of a test")
	}
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		panic(fmt.Sprintf("Failed to create Request; err=%s", err))
	}
	return req
}

// ! Call only inside test functions !
//
// [MustRequest] should be called as a wrapper for [http.Response] results.
// Performs basic checking with expected values.
func ResponseCheck(t *testing.T, ts *httptest.Server, tc *http.Client, subTestName string, expectedStatusCode int, resp *http.Response, err error) {
	if !testing.Testing() {
		panic("ResponseCheck called outside of a test")
	}
	if err != nil {
		logger.Log.Error("(%s) failed; err=%s", subTestName, err)
		t.Fatal(err)
	}
	if resp.StatusCode != expectedStatusCode {
		logger.Log.Error("(%s) failed; expected=%d; statuscode=%d", subTestName, expectedStatusCode, resp.StatusCode)
		t.Fatal("Invalid status code")
	}
}

// Basic server initalization and pinging "/" url
func TestPing(t *testing.T) {
	InitBackend(t)
	defer db.Conn.Close()
	defer logger.Log.Info("[TESTING FINISHED]")

	ts := httptest.NewServer(MuxHTTP())
	defer ts.Close()

	if _, err := http.Get(ts.URL); err != nil {
		t.Error(err)
		logger.Log.Error("req=GET failed; error=%s", err)
	}
}
