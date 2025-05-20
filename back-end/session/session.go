/*
Implements:
  - session can be encrypted into jwt for transfering
  - contain functionality to work with expiration time
  - contain functionality to key changing to secure the encryption

Package can be used in multithreads.
*/
package session

import "time"

//go:generate go tool github.com/princjef/gomarkdoc/cmd/gomarkdoc -o documentation.md

// Session expiration time messarued in hours
const EXPIRATION_TIME = 1.0

/*
Stores all info about session.

All fields can be accessed outside of a package.
The methods only here as helpers implementation.
*/
type Session struct {
	Name   string    `json:"name"`
	Expire time.Time `json:"expire"`
}

func New(name string) Session {
	return Session{name, time.Now().Add(EXPIRATION_TIME * time.Hour)}
}

func (ses *Session) UpdateExpiration() {
	ses.Expire = time.Now().Add(EXPIRATION_TIME * time.Hour)
}

func (ses Session) IsExpired() bool {
	return ses.Expire.Before(time.Now())
}
