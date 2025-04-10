package judge

import (
	"errors"
	"strings"

	"github.com/TrueHopolok/braincode-/judge/bf"
)

type bfGenerator bf.ByteCode

// NewBFGenerator creates a new brainfunk input generator.
//
// # Expected brainfunk API
//
// There will be no input provided.
// First output byte will be used as a group delimiter, second as a test delimiter. Delimiters cannot be the same.
// After that one or more groups of test may be provided.
//
// Tests must be separated by test delimiter. Groups must be separated by group delimiter.
// No terminating delimiter is required if last test is not empty.
//
// There is a hard cap of 1 billion steps and 64MiB memory. These may change in future versions.
func NewBFGenerator(source string) (InputGenerator, error) {
	bc, err := bf.Compile(source, -1)
	return &smartGenerator{bfGenerator(bc)}, err
}

func (b bfGenerator) GenerateInput() ([][]string, error) {
	var w bfGenWriter
	s := bf.NewState(bf.ByteCode(b), nil, &w, 1e9, 64e6)
	if err := s.Run(); err != nil {
		return nil, err
	}
	if err := w.Close(); err != nil {
		return nil, err
	}

	return w.groups, nil
}

type bfGenWriter struct {
	off        uint
	groupDelim byte
	testDelim  byte

	groups [][]string
	group  []string
	test   strings.Builder
}

func (w *bfGenWriter) Write(data []byte) (int, error) {
	for i, b := range data {
		if err := w.WriteByte(b); err != nil {
			return i, err
		}
	}
	return len(data), nil
}

func (w *bfGenWriter) WriteByte(c byte) error {
	defer func() {
		w.off++
	}()

	if w.off == 0 {
		// this is a group delim
		w.groupDelim = c
		return nil
	}
	if w.off == 1 {
		// this is a test delim
		w.testDelim = c

		if w.testDelim == w.groupDelim {
			return errors.New("test delimiter and group delimiters must be different")
		}

		return nil
	}

	if c != w.testDelim && c != w.groupDelim {
		w.test.WriteByte(c)
		return nil
	}

	// flush test
	w.group = append(w.group, w.test.String())
	w.test = strings.Builder{}

	if c == w.groupDelim {
		// flush group
		w.groups = append(w.groups, w.group)
		w.group = nil
	}

	return nil
}

func (w *bfGenWriter) Close() error {
	if w.off <= 1 {
		return errors.New("test header not found (expected at least 2 bytes)")
	}

	if w.test.Len() > 0 {
		w.group = append(w.group, w.test.String())
		w.test = strings.Builder{}
	}

	if len(w.group) > 0 {
		w.groups = append(w.groups, w.group)
		w.group = nil
	}

	return nil
}
