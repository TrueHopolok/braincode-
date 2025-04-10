package lua

import (
	"bytes"
	"cmp"
	"context"
	"encoding/gob"
	"errors"
	"fmt"
	"strings"
	"time"

	lua "github.com/yuin/gopher-lua"
	"github.com/yuin/gopher-lua/ast"
	"github.com/yuin/gopher-lua/parse"
)

// Checker contains a parsed lua script to be used for checking solutions to a problem.
//
// Script is sandboxed and global state is wiped between tests.
// Standard output is be captured and used as additional checker comments.
//
// Following formats are accepted:
//
// # Solution mode
//
// A function with signature solution(input) must be defined. It must return a string.
// It will be called with appropriate test input. Test will fail if submition result
// differs from the string returned from solution. Great for problems with singular answer.
//
// # Checker mode
//
// A function with signature checker(input, output) must be defined.
// Test input and submition output will be passes to the function, and it must decide if the test fails.
//   - no return value, nil, empty string or false are considered a successful test.
//   - true value is considered a failure without comment.
//   - any other value is considered a comment, will be stringified using tostring() and passes to the judge.
//
// Checker mode is always preferred over solution mode. At least one of 2 functions must be defined.
//
// Zero value checker is invalid, use [NewChecker] to construct one. It is safe to copy and use concurrently, because it is immutable.
type Checker struct {
	parsed      []ast.Stmt
	compiled    *lua.FunctionProto
	useSolution bool
}

// NewChecker parses source and creates a new Checker.
func NewChecker(source string) (Checker, error) {
	chunks, err := parse.Parse(strings.NewReader(source), "checker.lua")
	if err != nil {
		return Checker{}, fmt.Errorf("parse failed: %w", err)
	}
	return newChecker(chunks)
}

func newChecker(chunks []ast.Stmt) (Checker, error) {
	f, err := lua.Compile(chunks, "checker.lua")
	if err != nil {
		return Checker{}, fmt.Errorf("compilation failed: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	s := newLuaState()
	defer s.Close()
	s.SetContext(ctx)

	lfunc := s.NewFunctionFromProto(f)
	s.Push(lfunc)
	if err := s.PCall(0, lua.MultRet, nil); err != nil {
		return Checker{}, fmt.Errorf("initial execution failed: %w", err)
	}

	sol := s.GetGlobal("solution") != lua.LNil
	che := s.GetGlobal("checker") != lua.LNil

	if !sol && !che {
		return Checker{}, errors.New("either checker or solution must be defined")
	}

	return Checker{
		parsed:      chunks,
		compiled:    f,
		useSolution: !che,
	}, nil
}

// CheckOutput runs the checker.
// Non nil error means that checker has failed.
// If error is nil, string can be examined for test result.
// If string is empty, test passes.
// Otherwise, string will contain a possibly multiline checker comment.
func (c Checker) CheckOutput(input, output string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	l := newLuaState()
	defer l.Close()
	l.SetContext(ctx)

	lfunc := l.NewFunctionFromProto(c.compiled)
	l.Push(lfunc)
	if err := l.PCall(0, lua.MultRet, nil); err != nil {
		return "", err
	}

	if c.useSolution {
		return c.runSolution(l, input, output)
	} else {
		return c.runChecker(l, input, output)
	}
}

func (c Checker) runSolution(l luaState, input, output string) (string, error) {
	if err := l.CallByParam(lua.P{
		Fn:      l.GetGlobal("solution"),
		NRet:    1,
		Protect: true,
	}, lua.LString(input)); err != nil {
		return "", err
	}

	ret := l.CheckString(-1)

	if ret != output {
		return l.Buffer.String() + "result do not match", nil
	}

	return "", nil
}

func (c Checker) runChecker(l luaState, input, output string) (string, error) {
	if err := l.CallByParam(lua.P{
		Fn:      l.GetGlobal("checker"),
		NRet:    1,
		Protect: true,
	}, lua.LString(input), lua.LString(output)); err != nil {
		return "", err
	}

	ret := l.Get(-1)

	if lua.LVAsBool(ret) && ret != lua.LString("") {
		if ret == lua.LTrue {
			return cmp.Or(l.Buffer.String(), "test failed"), nil
		}

		return l.Buffer.String() + ret.String(), nil
	}

	return "", nil
}

func (c Checker) AppendBinary(b []byte) ([]byte, error) {
	// only ast is serialized
	buf := bytes.NewBuffer(b)
	err := gob.NewEncoder(buf).Encode(c.parsed)
	return buf.Bytes(), err
}

func (c *Checker) MarshalBinary() (data []byte, err error) {
	return c.AppendBinary(nil)
}

func (c *Checker) UnmarshalBinary(data []byte) error {
	var p []ast.Stmt

	if err := gob.NewDecoder(bytes.NewReader(data)).Decode(&p); err != nil {
		return err
	}

	cc, err := newChecker(p)
	if err != nil {
		return err
	}

	*c = cc

	return nil
}
