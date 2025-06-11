package controllers

import (
	"crypto/rand"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/TrueHopolok/braincode-/server/logger"
	"golang.org/x/crypto/argon2"
)

const (
	PSH_TIME = 1
	PSH_MEM  = 64 * 1024
	PSH_THR  = 4
	PSH_LEN  = 64
)

// PSH - Password Salted then Hashed
//
// Function that get password, salt, combines them then hashes using Argon2 method.
// Result can be used in database for checking the similarities.
func PSH(pass string, salt []byte) []byte {
	return argon2.IDKey([]byte(pass), salt, PSH_TIME, PSH_MEM, PSH_THR, PSH_LEN)
}

// Generates crypto-random salt for a user's password.
func SaltGen() []byte {
	salt := make([]byte, PSH_LEN)
	if _, err := rand.Read(salt); err != nil {
		logger.Log.Fatal("crypto rand raised error=%s", err)
	}
	return salt
}

// Return 2 booleans that tell what language does user want.
// In case of invalid parameter, function will return false in isvalid field.
// On invalid parameter will output an error, thus this must be last write into response.
func langHandler(w http.ResponseWriter, r *http.Request) (isvalid bool, isenglish bool) {
	lang := strings.ToLower(r.URL.Query().Get("lang"))
	if lang == "ru" {
		logger.Log.Debug("req=%p lang=%s(RU) was selected", r, lang)
		return true, false
	} else if lang != "" && lang != "en" {
		// FIXME(anpir): this should be 404 / 400, not 406
		// Quoting MDN: https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Status/406
		//
		//    The HTTP 406 Not Acceptable client error response status code indicates that the server could not
		//    produce a response matching the list of acceptable values defined in the request's proactive content
		//    negotiation headers and that the server was unwilling to supply a default representation.
		//
		//    ...
		//
		//    If a server returns a 406, the body of the message should contain the list of available representations
		//    for the resource, allowing the user to choose, although no standard way for this is defined.
		//
		// There is no proactive content negotiation. No available representations are returned.
		http.Error(w, "Such language selection is not allowed", 406)
		logger.Log.Debug("req=%p lang=%s is not allowed", r, lang)
		return false, false
	}
	logger.Log.Debug("req=%p lang=%s(EN) was selected", r, lang)
	return true, true
}

// This should be the last write into the response!
//
// Redirects user to main/problemset page with 303 error code (SeeOther)
// Output in logger the given action for debugging purposes
func redirect2main(w http.ResponseWriter, r *http.Request, action string) {
	http.Redirect(w, r, "/", http.StatusSeeOther)
	logger.Log.Debug("req=%p user redirect after %s", r, action)
}

// This should be the last write into the response!
//
// Redirect user to their profile page with 303 error code (SeeOther)
// Output in logger the given action for debugging purposes
func redirect2stats(w http.ResponseWriter, r *http.Request, action string) {
	http.Redirect(w, r, "/stats/", 303)
	logger.Log.Debug("req=%p user redirect after %s", r, action)
}

func redirectError(w http.ResponseWriter, r *http.Request, errcode int) {
	redirectErrorString(w, r, strconv.Itoa(errcode))
}

func redirectErrorString(w http.ResponseWriter, r *http.Request, errorS string) {
	url := "?error=" + url.QueryEscape(errorS)
	if r.URL.Query().Has("lang") {
		url += "&lang=" + r.URL.Query().Get("lang")
	}
	http.Redirect(w, r, url, http.StatusSeeOther)
	logger.Log.Debug("req=%p user redirect with error %v", r, errorS)
}

// This should be the last write into the response!
//
// Can be used to temporarly substitue some code.
// Will return 503 error code (ServiceUnavailable) as http response and write into a logger that given feature is not implemented.
func errResp_NotImplemented(w http.ResponseWriter, r *http.Request, feature string) {
	feature = fmt.Sprintf("%s not implemented yet", feature)
	http.Error(w, feature, http.StatusServiceUnavailable)
	logger.Log.Error("req=%p failed; error=%s", r, feature)
}

// This should be the last write into the response!
//
// Write an error into the both response and logger.
// Should be used in case some internal error in execution happened,
// not for user invalid request handling.
func errResp_Fatal(w http.ResponseWriter, r *http.Request, err error) {
	http.Error(w, "Failed to write into the response body", http.StatusInternalServerError)
	logger.Log.Error("req=%p failed; error=%s", r, err)
}

// This should be the last write into the response!
//
// Write into response that given content-type is not allowed.
// Also writes into the logger for debbuging purposes.
func denyResp_ContentTypeNotAllowed(w http.ResponseWriter, r *http.Request, allowedTypes ...string) {
	result := fmt.Sprintf("Content-Type=%s is not allowed\nAllowed=", r.Header.Get("Content-Type"))
	for _, contenttype := range allowedTypes {
		result += contenttype
	}
	http.Error(w, result, http.StatusNotAcceptable)
	logger.Log.Debug("req=%p Content-Type=%s is not allowed", r, r.Header.Get("Content-Type"))
}
