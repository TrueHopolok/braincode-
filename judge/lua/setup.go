// Package lua implements [judge.InputGenerator] and [judge.OutputChecker] for lua programs.
package lua

//go:generate go tool github.com/princjef/gomarkdoc/cmd/gomarkdoc -o documentation.md

import (
	"bytes"
	"fmt"
	"math/rand/v2"
	"strings"
	"sync"

	lua "github.com/yuin/gopher-lua"
)

type luaState struct {
	*lua.LState
	*bytes.Buffer // standard output
}

func newLuaState() luaState {
	b := new(bytes.Buffer)

	l := lua.NewState(lua.Options{
		CallStackSize:   256,
		RegistryMaxSize: 1024 * 128,
		SkipOpenLibs:    true,
	})

	lua.OpenBase(l)
	lua.OpenMath(l)
	lua.OpenString(l)
	lua.OpenTable(l)

	// not part of standard in the first place
	deleteFunction(l, "_printregs")
	deleteFunction(l, "newproxy")

	// can access filesystem
	denyFunction(l, "dofile")
	denyFunction(l, "loadfile")

	rng := rand.New(rand.NewPCG(123, 456))
	var rngMu sync.Mutex

	// these are overriden to make the execution deterministic

	l.SetGlobal("random", l.NewFunction(func(l *lua.LState) int {
		rngMu.Lock()
		defer rngMu.Unlock()

		switch l.GetTop() {
		case 0:
			l.Push(lua.LNumber(rng.Float64()))
		case 1:
			n := l.CheckInt(1)
			l.Push(lua.LNumber(rng.IntN(n) + 1))
		default:
			min := l.CheckInt(1)
			max := l.CheckInt(2) + 1
			l.Push(lua.LNumber(rng.IntN(max-min) + min))
		}
		return 1

	}))

	l.SetGlobal("randomseed", l.NewFunction(func(l *lua.LState) int {
		rngMu.Lock()
		defer rngMu.Unlock()
		rng = rand.New(rand.NewPCG(uint64(l.CheckInt64(1)), 42))
		return 0
	}))

	l.SetGlobal("print", l.NewFunction(func(l *lua.LState) int {
		var s []string
		for i := 1; i <= l.GetTop(); i++ {
			s = append(s, l.Get(i).String())
		}
		fmt.Fprintf(b, "%s\n", strings.Join(s, " "))
		return 0
	}))

	return luaState{l, b}
}

func denyFunction(l *lua.LState, name string) {
	l.SetGlobal(name, l.NewFunction(func(l *lua.LState) int {
		l.RaiseError("function %s is not available", name)
		return 0
	}))
}

func deleteFunction(l *lua.LState, name string) {
	l.SetGlobal(name, lua.LNil)
}
