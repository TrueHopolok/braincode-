package controllers

import (
	"net/http"

	"github.com/TrueHopolok/braincode-/server/logger"
	"github.com/TrueHopolok/braincode-/server/session"
)

// Function is used for most cases except: registration, login and logout requests.
// Get session token from the request header.
// Validates it.
// If valid: updates the token and writes it into a header.
func sessionHandler(w http.ResponseWriter, r *http.Request) (session.Session, bool) {
	token := r.Header.Get("session")
	var ses session.Session
	isauth := ses.ValidateJWT(token) && !ses.IsExpired()
	if isauth {
		ses.UpdateExpiration()
		w.Header().Add("session", ses.CreateJWT())
	}
	return ses, isauth
}

// Can be used to temporarly substitue some code.
// Will return 503 error code as http response and write into a logger an given message.
func notImplemented(w http.ResponseWriter, r *http.Request, e string) {
	http.Error(w, "Not implemented yet", 503)
	logger.Log.Error("req=%p failed; error=%s", r, e)
}
