package controllers

import (
	"crypto/rand"
	"fmt"
	"net/http"

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
func langHandler(w http.ResponseWriter, r *http.Request) (isenglish bool, isvalid bool) {
	lang := r.URL.Query().Get("lang")
	if lang == "ru" {
		return false, true
	} else if lang != "" && lang != "en" {
		http.Error(w, "Such language selection is not allowed", 406)
		logger.Log.Debug("req=%p lang=%s is not allowed", r, lang)
		return false, false
	}
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
// Used to display not allowed method error.
// Will add all provided methods into the both response and logger.
func denyResp_MethodNotAllowed(w http.ResponseWriter, r *http.Request, allowedMethods ...string) {
	result := fmt.Sprintf("Method=%s is not allowed\nAllowed=", r.Method)
	for _, method := range allowedMethods {
		result += method
		w.Header().Add("Allow", method)
	}
	http.Error(w, result, http.StatusMethodNotAllowed)
	logger.Log.Debug("req=%p method=%s is not allowed", r, r.Method)
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
