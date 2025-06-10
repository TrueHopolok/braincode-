package main

import (
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
	t.Helper()
	if !testing.Testing() {
		panic("InitBackend called outside of a test")
	}

	// TODO(anpir): use other DB for tests (braincode-testing)
	// which is wiped on start each tests.
	// Would have been a lot easier, if migrations were 2-way :/

	//* Config init for testing
	config.OverrideConfig(t, config.Config{
		Verbose:       true,
		EnableConsole: false,
		LogFilepath:   "server.log",
		TemplatesPath: "../frontend/",
		StaticPath:    "../frontend/static/",
		DBuser:        "root",
		DBpass:        "root",
		DBname:        "braincode_test",
		DBqueriesPath: "db/queries/",
		Secure:        false,
	})

	//* Logger init
	logger.Testing()
	logger.Log.Info("[TESTING STARTED]")

	//* Database init
	db.InitTesting(t)

	//* Database migrate
	if err := db.Migrate(); err != nil {
		t.Fatalf("database migration failed: err = %v", err)
		logger.Log.Error("Migration: execution failed; error=%s", err)
		db.Conn.Close()
		logger.Log.Info("[TESTING FINISHED]")
		return
	}

	//* Templates init
	if err := prepared.Init(); err != nil {
		t.Fatalf("template initialization failed: err = %v", err)
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
	t.Helper()
	if !testing.Testing() {
		panic("MustRequest called outside of a test")
	}
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		t.Fatalf("failed to create a request: err=%v", err)
	}
	return req
}

// ! Call only inside test functions !
//
// [MustRequest] should be called as a wrapper for [http.Response] results.
// Performs basic checking with expected values.
func ResponseCheck(t *testing.T, ts *httptest.Server, tc *http.Client, subTestName string, expectedStatusCode int, resp *http.Response, err error) {
	t.Helper()
	if !testing.Testing() {
		panic("ResponseCheck called outside of a test")
	}
	if err != nil {
		logger.Log.Error("(%s) failed; err=%s", subTestName, err)
		t.Fatalf("request failed err = %v", err)
	}
	if resp.StatusCode != expectedStatusCode {
		logger.Log.Error("(%s) failed; expected=%d; status=%s", subTestName, expectedStatusCode, resp.Status)
		t.Fatalf("bad status: status = %s; want %d", resp.Status, expectedStatusCode)
	}

	var cc []*http.Cookie
	for _, s := range resp.Header["Set-Cookie"] {
		c, err := http.ParseSetCookie(s)
		if err != nil {
			t.Errorf("server set invalid cookie %q: %v", s, err)
		}
		cc = append(cc, c)
	}
	if len(cc) > 0 {
		tc.Jar.SetCookies(resp.Request.URL, cc)
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
