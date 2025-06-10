package ml_test

import (
	"strings"
	"testing"

	"github.com/TrueHopolok/braincode-/judge/ml"
	"github.com/TrueHopolok/braincode-/judge/ml/testhelper"
)

func TestDocument_Serialization(t *testing.T) {
	data := `
	.task = My task
	.section = Hello!
	.paragraph
	Hello, this is a paragraph.
	..
	`

	was, err := ml.Parse(strings.NewReader(data))
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}

	binary, err := was.MarshalBinary()
	if err != nil {
		t.Fatalf("serialization failed: %v", err)
	}

	var got ml.Document
	if err := got.UnmarshalBinary(binary); err != nil {
		t.Fatalf("deserialization failed: %v", err)
	}

	testhelper.DiffValues(t, was, got)
}
