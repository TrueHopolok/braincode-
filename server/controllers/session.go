package controllers

import (
	"net/http"

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
