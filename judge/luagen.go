package judge

import "github.com/TrueHopolok/braincode-/judge/lua"

type luaGenerator string

// NewLuaGenerator create a new generator from lua source code.
// See [lua.GetTests] for details.
func NewLuaGenerator(source string) InputGenerator {
	return &smartGenerator{luaGenerator(source)}
}

func (g luaGenerator) GenerateInput() ([][]string, error) {
	return lua.GetTests(string(g))
}
