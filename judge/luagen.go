package judge

import "github.com/TrueHopolok/braincode-/judge/lua"

type luaGenerator string

func NewLuaGenerator(source string) InputGenerator {
	return &smartGenerator{luaGenerator(source)}
}

func (g luaGenerator) GenerateInput() ([][]string, error) {
	return lua.GetTests(string(g))
}
