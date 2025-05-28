package controllers

import (
	"crypto/rand"
	"fmt"
	"net/http"

	"github.com/TrueHopolok/braincode-/server/logger"
	"github.com/TrueHopolok/braincode-/server/session"
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

// Function is used for most cases except: registration, login and logout requests.
// Get session token from the request header.
// Validates it.
// If valid: updates the token and writes it into a header.
// If invalid or expired: return empty string as a name and false in aiauth field.
func sessionHandler(w http.ResponseWriter, r *http.Request) (string, bool) {
	token := r.Header.Get("Session")
	var ses session.Session
	isauth := ses.ValidateJWT(token) && !ses.IsExpired()
	if isauth {
		ses.UpdateExpiration()
		w.Header().Add("Session", ses.CreateJWT())
	} else {
		ses.Name = ""
	}
	return ses.Name, isauth
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
	http.Redirect(w, r, "/", 303)
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
// Will return 503 error code as http response and write into a logger that given feature is not implemented.
func errResp_NotImplemented(w http.ResponseWriter, r *http.Request, feature string) {
	feature = fmt.Sprintf("%s not implemented yet", feature)
	http.Error(w, feature, 503)
	logger.Log.Error("req=%p failed; error=%s", r, feature)
}

// This should be the last write into the response!
//
// Write an error into the both response and logger.
// Should be used in case some internal error in execution happened,
// not for user invalid request handling.
func errResp_Fatal(w http.ResponseWriter, r *http.Request, err error) {
	http.Error(w, "Failed to write into the response body", 500)
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
	http.Error(w, result, 405)
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
	http.Error(w, result, 406)
	logger.Log.Debug("req=%p Content-Type=%s is not allowed", r, r.Header.Get("Content-Type"))
}

// This should be the last write into the response!
//
// Output that user is not authorized to use this page.
// Writes in both logger and response.
func denyResp_NotAuthorized(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Login into account to use this page", 401)
	logger.Log.Debug("req=%p user trying to access page while unauthorized", r)
}
