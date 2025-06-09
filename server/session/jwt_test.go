package session

import (
	"math/rand/v2"
	"strings"
	"testing"
	"time"
)

func ReadDeterministic(seed1, seed2 uint64, buf []byte) {
	rng := rand.New(rand.NewPCG(seed1, seed2))
	for i := range buf {
		buf[i] = byte(rng.Uint32())
	}
}

func TestSession_ValidateJWT_correct(t *testing.T) {
	keys.cur = make([]byte, key_size)
	keys.prv = make([]byte, key_size)
	ReadDeterministic(4, 2, keys.cur)
	ReadDeterministic(2, 4, keys.prv)
	s := Session{
		Name:   "hello world!",
		Expire: time.Date(2020, 11, 11, 11, 11, 11, 111, time.UTC),
	}

	jwt := s.CreateJWT()

	var s2 Session

	if !s2.ValidateJWT(jwt) {
		t.Fatal("invalid jwt")
	}

	if s != s2 {
		t.Errorf("session do not match: before %+v, after %+v", s, s2)
	}
}

func TestSession_ValidateJWT_incorrect(t *testing.T) {
	keys.cur = make([]byte, key_size)
	keys.prv = make([]byte, key_size)
	ReadDeterministic(4, 2, keys.cur)
	ReadDeterministic(2, 4, keys.prv)

	cases := [][2]string{
		{"other jwt", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWUsImlhdCI6MTUxNjIzOTAyMn0.KMUFsIDTnFmyG3nMiGM6H9FNFUROf3wh7SmqJp-QV30"},
		{"empty string", ""},
		{"random text", "invalid bla bla bla bla"},
		{"a lot of dots", strings.Repeat(".a", 1000000)},
	}

	for _, tc := range cases {
		t.Run(tc[0], func(t *testing.T) {
			t.Parallel()
			var s Session
			if s.ValidateJWT(tc[1]) {
				t.Error("token reported as valid")
			}

		})
	}
}
