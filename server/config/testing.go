package config

import (
	"sync"
	"testing"
)

func OverrideConfig(t *testing.T, cfg Config) {
	// require called to pass in testing.T to discourage calling outside of tests.
	// crash the program if caller forged t.
	if !testing.Testing() {
		panic("OverrideConfig called outside of a test")
	}
	once = sync.OnceValue(func() Config {
		return cfg
	})
}
