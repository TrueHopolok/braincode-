package session

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"strings"
)

const header = "{'alg':'HS256','typ':'JWT'}" // Standard JWT header
var b64Header = base64.URLEncoding.EncodeToString([]byte(header))

func tokenize(header, body string) string {
	hash := hmac.New(sha256.New, keys.cur)
	_, err := hash.Write([]byte(header + body))
	if err != nil {
		panic(err)
	}
	return base64.URLEncoding.EncodeToString(hash.Sum(nil))
}

// Create jwt with information from provided session.
//
// May panic if somehow JSON serialization or HMAC hash fail.
func (ses Session) CreateJWT() string {
	jsonBody, err := json.Marshal(&ses)
	if err != nil {
		panic(err)
	}
	b64Body := base64.URLEncoding.EncodeToString(jsonBody)
	return b64Header + "." + b64Body + "." + tokenize(b64Header, b64Body)
}

// ValidateJWT parses given token into ses, reporting whether token is valid.
//
// May panic if token was create on a different version or given a different JSON schema.
//
// Receiver is unchanged in case token is invalid.
//
// This function is safe to use with user input.
func (ses *Session) ValidateJWT(token string) bool {
	fields := strings.SplitN(token, ".", 4)
	if len(fields) != 3 {
		return false
	}
	if fields[0]+fields[1]+tokenize(fields[0], fields[1]) != token {
		return false
	}
	if err := json.Unmarshal([]byte(fields[1]), ses); err != nil {
		// TODO(anpir): this may panic if structure of token changes;
		// probably should return false instead.
		panic(err)
	}
	return true
}
