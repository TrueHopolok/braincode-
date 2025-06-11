package main

import (
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"net/url"
	"testing"

	db "github.com/TrueHopolok/braincode-/server/db"
	"github.com/TrueHopolok/braincode-/server/logger"
	"golang.org/x/net/publicsuffix"
)

// Test server all authefication processes in this order:
//   - Submit solution (fail=notauth),
//   - Register (ok),
//   - Register (fail=isauth),
//   - Logout (ok),
//   - Logout (fail=notauth),
//   - Register (fail=already_exists),
//   - Login (fail=invalid_username),
//   - Login (fail=invalid_password),
//   - Login (fail=invalid_both_authdata),
//   - Login (ok),
//   - Login (fail=isauth),
//   - Delete user (ok),
//   - Delete user (fail=notauth)
//   - Login (fail=dont_exists);
//
// Auth data:
//   - Username: "Tester",
//   - Password: "Password";
func TestAuth(t *testing.T) {
	InitBackend(t)
	defer db.Conn.Close()
	defer logger.Log.Info("[TESTING FINISHED]")

	ts := httptest.NewServer(MuxHTTP())
	defer ts.Close()
	tc := ts.Client()
	tc.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}
	var err error
	tc.Jar, err = cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	if err != nil {
		t.Fatalf("cookie jar init failed: err = %v", err)
	}
	var req *http.Request
	// var ses session.Session
	var resp *http.Response
	var subTestName string
	var expectedStatusCode int

	// //* Submit solution (fail=notauth)
	// subTestName = "Submit solution (fail=notauth)"
	// expectedStatusCode = http.StatusUnauthorized
	// resp, err = tc.Post(ts.URL+"/task/", "application/brainfunk", nil) // [body] is not necessary since request should fail with http.StatusUnauthorized (401)
	// ResponseCheck(t, ts, tc, subTestName, expectedStatusCode, resp, err)

	//* Register (ok)
	subTestName = "Register (ok)"
	expectedStatusCode = http.StatusSeeOther
	resp, err = tc.PostForm(ts.URL+"/register/", url.Values{"username": {"Tester"}, "password": {"Password"}})
	ResponseCheck(t, ts, tc, subTestName, expectedStatusCode, resp, err)

	// //* Register (fail=isauth)
	// subTestName = "Register (fail=isauth)"
	// expectedStatusCode = http.StatusBadRequest
	// req = MustRequest(t, "DELETE", ts.URL+"/register/", nil)
	// resp, err = tc.PostForm(ts.URL+"/register/", url.Values{"username": {"Tester"}, "password": {"Password"}})
	// ResponseCheck(t, ts, tc, subTestName, expectedStatusCode, resp, err)

	//* Logout (ok)
	subTestName = "Logout (ok)"
	expectedStatusCode = http.StatusSeeOther
	req = MustRequest(t, "POST", ts.URL+"/logout/", nil)
	resp, err = tc.Do(req)
	// This test keeps failing, i dont know why, auth works fine in prod
	ResponseCheck(t, ts, tc, subTestName, expectedStatusCode, resp, err)

	// //* Logout (fail=notauth)
	// subTestName = "Logout (fail=notauth)"
	// expectedStatusCode = http.StatusUnauthorized
	// resp, err = tc.Do(MustRequest(t, "DELETE", ts.URL+"/login/", nil))
	// ResponseCheck(t, ts, tc, subTestName, expectedStatusCode, resp, err)

	// //* Register (fail=already_exists)
	// subTestName = "Register (fail=already_exists)"
	// expectedStatusCode = http.StatusNotAcceptable
	// resp, err = tc.PostForm(ts.URL+"/register/", url.Values{"username": {"Tester"}, "password": {"Password"}})
	// ResponseCheck(t, ts, tc, subTestName, expectedStatusCode, resp, err)

	// //* Login (fail=invalid_username)
	// subTestName = "Login (fail=invalid_username)"
	// expectedStatusCode = http.StatusNotAcceptable
	// resp, err = tc.PostForm(ts.URL+"/register/", url.Values{"username": {"NotTester"}, "password": {"Password"}})
	// ResponseCheck(t, ts, tc, subTestName, expectedStatusCode, resp, err)

	// //* Login (fail=invalid_password)
	// subTestName = "Login (fail=invalid_username)"
	// expectedStatusCode = http.StatusNotAcceptable
	// resp, err = tc.PostForm(ts.URL+"/register/", url.Values{"username": {"Tester"}, "password": {"Qwerty123"}})
	// ResponseCheck(t, ts, tc, subTestName, expectedStatusCode, resp, err)

	// //* Login (fail=invalid_both_authdata)
	// subTestName = "Login (fail=invalid_username)"
	// expectedStatusCode = http.StatusNotAcceptable
	// resp, err = tc.PostForm(ts.URL+"/register/", url.Values{"username": {"NotTester"}, "password": {"Qwerty123"}})
	// ResponseCheck(t, ts, tc, subTestName, expectedStatusCode, resp, err)

	//* Login (ok)
	subTestName = "Login (ok)"
	expectedStatusCode = http.StatusSeeOther
	resp, err = tc.PostForm(ts.URL+"/login/", url.Values{"username": {"Tester"}, "password": {"Password"}})
	ResponseCheck(t, ts, tc, subTestName, expectedStatusCode, resp, err)

	// //* Login (fail=isauth)
	// subTestName = "Login (fail=isauth)"
	// expectedStatusCode = http.StatusBadRequest
	// resp, err = tc.PostForm(ts.URL+"/login/", url.Values{"username": {"Tester"}, "password": {"Password"}})
	// ResponseCheck(t, ts, tc, subTestName, expectedStatusCode, resp, err)

	//* Delete user (ok)
	subTestName = "Delete user (ok)"
	expectedStatusCode = http.StatusSeeOther
	req = MustRequest(t, "DELETE", ts.URL+"/stats/", nil)
	resp, err = tc.Do(req)
	ResponseCheck(t, ts, tc, subTestName, expectedStatusCode, resp, err)

	// //* Delete user (fail=notauth)
	// subTestName = "Delete user (fail=notauth)"
	// expectedStatusCode = http.StatusUnauthorized
	// resp, err = tc.Do(MustRequest(t, "DELETE", ts.URL+"/stats/", nil))
	// ResponseCheck(t, ts, tc, subTestName, expectedStatusCode, resp, err)

	// //* Login (fail=dont_exists)
	// subTestName = "Login (fail=dont_exists)"
	// expectedStatusCode = http.StatusNotAcceptable
	// resp, err = tc.PostForm(ts.URL+"/login/", url.Values{"username": {"Tester"}, "password": {"Password"}})
	// ResponseCheck(t, ts, tc, subTestName, expectedStatusCode, resp, err)
}
