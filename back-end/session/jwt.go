package session

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"strings"
)

const (
	header = "{'alg':'HS256','typ':'JWT'}" // Standart JWT header
	b64_header = "eydhbGcnOidIUzI1NicsJ3R5cCc6J0pXVCd9" // Generated via base64.URLEncoding.EncodeToString([]byte(header))
)

func tokenize(header, body string) string {
	hash := hmac.New(sha256.New, keys.cur)
	_, err := hash.Write([]byte(header+body))
	if err != nil {
		panic(err)
	}
	return base64.URLEncoding.EncodeToString(hash.Sum(nil))
}

/*
Create jwt with information from provided session.
May panic if somehow json or hmac packages fail.
*/
func (ses Session) CreateJWT() string {
	js_body, err := json.Marshal(ses)
	if err != nil {
		panic(err)
	}
	b64_body := base64.URLEncoding.EncodeToString(js_body)
	return b64_header + "." + b64_body + "." + tokenize(b64_header, b64_body)
}


/*
Check if given jwt is valid using keysin memory.
Decoded information will save into provided session.
May panic if somehow json or hmac packages fail, but that should be possible.
	if jwt is invalid: ses is unchanged
*/
func (ses *Session) ValidateJWT(token string) bool {
	fields := strings.Split(token, ".")
	if len(fields) != 3 {
		return false
	}
	if fields[0] + fields[1] + tokenize(fields[0], fields[1]) != token {
		return false
	}
	if err := json.Unmarshal([]byte(fields[1]), ses); err != nil {
		panic(err)
	}
	return true
}
