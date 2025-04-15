package session

import (
	"bytes"
	"encoding/gob"
)

// Encrypts session into slice of bytes
func (ses Session) Cypher() (raw []byte, err error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err = enc.Encode(ses)
	if err != nil {
		return nil, err
	}
	raw = buf.Next(buf.Len())
	for i := 0; i < len(raw); i++ {
		raw[i] ^= keys.cur[i%key_size]
	}
	return raw, nil
}


/*
Decrypt slice of bytes into session

Return error:
	if raw is not Cypher(ses)
	if raw use too old encrpytion key
*/
func Decypher(raw []byte) (ses Session, err error) {
	for i := 0; i < len(raw); i++ {
		raw[i] ^= keys.cur[i%key_size]
	}
	var buf bytes.Buffer
	buf.Write(raw)
	dec := gob.NewDecoder(&buf)
	err = dec.Decode(&ses)
	if err != nil {
		for i := 0; i < len(raw); i++ {
			raw[i] ^= keys.cur[i%key_size] ^ keys.prv[i%key_size]
		}
		buf.Reset()
		buf.Write(raw)
		err = dec.Decode(&ses)
	}
	return ses, err
}
