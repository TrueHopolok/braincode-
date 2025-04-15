/*
Implements:
  - main functionality of sessions ... nothing for now;
  - encrytion and decryption.

Package can be used in multithreads.
*/
package session

//go:generate go tool github.com/princjef/gomarkdoc/cmd/gomarkdoc -o documentation.md

/*
Stores all info about session.

All fields can be accessed outside of a package.
The methods only here as helpers implementation.
*/
type Session struct {
	UserID int64
}

// Potentially extra functionality that session is neede for