package models

import (
	"bytes"
	"encoding/gob"
	"errors"

	"github.com/TrueHopolok/braincode-/judge"
)

// Function used to convert a problem struct provided by judge
// Into a slice of bytes in gob encoding
func marshalProblem(prb judge.Problem) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)

	byt, err := judge.MarshalGenerator(prb.InputGenerator)
	if err != nil {
		return nil, err
	}
	err = enc.Encode(byt)
	if err != nil {
		return nil, err
	}

	byt, err = judge.MarshalChecker(prb.OutputChecker)
	if err != nil {
		return nil, err
	}
	err = enc.Encode(byt)
	if err != nil {
		return nil, err
	}

	err = enc.Encode(prb.Instructions)
	if err != nil {
		return nil, err
	}

	err = enc.Encode(prb.Steps)
	if err != nil {
		return nil, err
	}

	err = enc.Encode(prb.Memory)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// Function used to extract all fields of given problem
// Into a struct provided by judge package
func unmarshalProblem(raw []byte) (judge.Problem, error) {
	var buf bytes.Buffer
	n, err := buf.Read(raw)
	if err != nil {
		return judge.Problem{}, err
	} else if n != len(raw) {
		return judge.Problem{}, errors.New("unmarshaling a problem, buffer reading failure happend")
	}
	dec := gob.NewDecoder(&buf)
	var prb judge.Problem
	var byt []byte

	err = dec.Decode(&byt)
	if err != nil {
		return judge.Problem{}, err
	}
	prb.InputGenerator, err = judge.UnmarshalGenerator(byt)
	if err != nil {
		return judge.Problem{}, err
	}

	err = dec.Decode(&byt)
	if err != nil {
		return judge.Problem{}, err
	}
	prb.OutputChecker, err = judge.UnmarshalChecker(byt)

	err = dec.Decode(&prb.Instructions)
	if err != nil {
		return judge.Problem{}, err
	}

	err = dec.Decode(&prb.Steps)
	if err != nil {
		return judge.Problem{}, err
	}

	err = dec.Decode(&prb.Memory)
	if err != nil {
		return judge.Problem{}, err
	}

	return prb, nil
}
