package judge

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"

	"github.com/TrueHopolok/braincode-/judge/lua"
)

type serializedGenerator struct {
	List *listGenerator
	BF   *bfGenerator
	Lua  *luaGenerator
}

func MarshalGenerator(g InputGenerator) ([]byte, error) {
	return AppendGenerator(g, nil)
}

func AppendGenerator(g InputGenerator, b []byte) ([]byte, error) {
	var val serializedGenerator
	switch gen := g.(type) {
	case bfGenerator:
		val.BF = &gen
	case listGenerator:
		val.List = &gen
	case luaGenerator:
		val.Lua = &gen
	case *smartGenerator:
		return AppendGenerator(gen.InputGenerator, b)
	default:
		return nil, fmt.Errorf("unexpected generator: %T", g)
	}

	buf := bytes.NewBuffer(b)
	err := gob.NewEncoder(buf).Encode(&val)
	return buf.Bytes(), err
}

func UnmarshalGenerator(b []byte) (InputGenerator, error) {
	var val serializedGenerator
	if err := gob.NewDecoder(bytes.NewReader(b)).Decode(&val); err != nil {
		return nil, err
	}

	if val.Lua != nil {
		return &smartGenerator{*val.Lua}, nil
	}
	if val.List != nil {
		return *val.List, nil
	}
	if val.BF != nil {
		return &smartGenerator{*val.BF}, nil
	}
	return nil, errors.New("value did not contain any known generator")
}

type serializedChecker struct {
	List       *listSolution
	BFSolution *bfSolution
	BFChecker  *bfChecker
	Lua        *lua.Checker
}

func MarshalChecker(c OutputChecker) ([]byte, error) {
	return AppendChecker(c, nil)
}

func AppendChecker(c OutputChecker, b []byte) ([]byte, error) {
	var val serializedChecker
	switch che := c.(type) {
	case bfChecker:
		val.BFChecker = &che
	case bfSolution:
		val.BFSolution = &che
	case listSolution:
		val.List = &che
	case luaChecker:
		val.Lua = che.Checker
	default:
		return nil, fmt.Errorf("unexpected checker: %T", che)
	}

	buf := bytes.NewBuffer(b)
	err := gob.NewEncoder(buf).Encode(&val)
	return buf.Bytes(), err
}

func UnmarshalChecker(b []byte) (OutputChecker, error) {
	var val serializedChecker
	if err := gob.NewDecoder(bytes.NewReader(b)).Decode(&val); err != nil {
		return nil, err
	}

	if val.Lua != nil {
		return luaChecker{val.Lua}, nil
	}
	if val.BFChecker != nil {
		return *val.BFChecker, nil
	}
	if val.BFSolution != nil {
		return *val.BFSolution, nil
	}
	if val.List != nil {
		return *val.List, nil
	}
	return nil, errors.New("value did not contain any known checker")
}

func AppendProblem(p Problem, b []byte) ([]byte, error) {
	panic("TODO")
}

func MarshalProblem(p Problem) ([]byte, error) {
	panic("TODO")
}

func UnmarshalProblem(b []byte) (Problem, error) {
	panic("TODO")
}
