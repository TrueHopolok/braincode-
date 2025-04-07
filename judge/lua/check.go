package lua

import (
	"bytes"
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

type Checker struct {
	parsed      []ast.Stmt
	compiled    *lua.FunctionProto
	useSolution bool
}

func NewChecker(source string) (*Checker, error) {
	chunks, err := parse.Parse(strings.NewReader(source), "checker.lua")
	if err != nil {
		return nil, fmt.Errorf("parse failed: %w", err)
	}
	return newChecker(chunks)
}

func newChecker(chunks []ast.Stmt) (*Checker, error) {

	f, err := lua.Compile(chunks, "checker.lua")
	if err != nil {
		return nil, fmt.Errorf("compilation failed: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	s := newLuaState()
	defer s.Close()
	s.SetContext(ctx)

	lfunc := s.NewFunctionFromProto(f)
	s.Push(lfunc)
	if err := s.PCall(0, lua.MultRet, nil); err != nil {
		return nil, fmt.Errorf("initial execution failed: %w", err)
	}

	sol := s.GetGlobal("solution") != lua.LNil
	che := s.GetGlobal("checker") != lua.LNil

	if !sol && !che {
		return nil, errors.New("either checker or solution must be defined")
	}

	return &Checker{
		parsed:      chunks,
		compiled:    f,
		useSolution: !che,
	}, nil
}

func (c *Checker) CheckOutput(input, output string) (string, error) {
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

func (c *Checker) runSolution(l luaState, input, output string) (string, error) {
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

func (c *Checker) runChecker(l luaState, input, output string) (string, error) {
	if err := l.CallByParam(lua.P{
		Fn:      l.GetGlobal("checker"),
		NRet:    1,
		Protect: true,
	}, lua.LString(input), lua.LString(output)); err != nil {
		return "", err
	}

	ret := l.Get(-1)

	if lua.LVAsBool(ret) {
		return l.Buffer.String() + ret.String(), nil
	}

	return "", nil
}

func (c *Checker) AppendBinary(b []byte) ([]byte, error) {
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

	*c = *cc

	return nil
}
