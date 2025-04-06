package lua

import (
	"math/rand/v2"
	"sync"

	lua "github.com/yuin/gopher-lua"
)

func newLuaState() *lua.LState {
	l := lua.NewState(lua.Options{
		CallStackSize:   256,
		RegistryMaxSize: 1024 * 128,
		SkipOpenLibs:    true,
	})

	lua.OpenBase(l)
	lua.OpenMath(l)
	lua.OpenString(l)
	lua.OpenTable(l)

	// will write to stdout
	noopFunction(l, "print")

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

	return l
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

func noopFunction(l *lua.LState, name string) {
	l.SetGlobal(name, l.NewFunction(func(l *lua.LState) int { return 0 }))
}
