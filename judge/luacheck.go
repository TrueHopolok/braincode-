package judge

import "github.com/TrueHopolok/braincode-/judge/lua"

type luaChecker struct{ *lua.Checker }

func NewLuaChecker(source string) (OutputChecker, error) {
	c, err := lua.NewChecker(source)
	return luaChecker{c}, err
}

func (l luaChecker) CheckOutput(input string, output string) Verdict {
	res, err := l.Checker.CheckOutput(input, output)
	if err != nil {
		return Verdict{
			Status:  StatusCheckerFailed,
			Comment: err.Error(),
		}
	}

	if res != "" {
		return Verdict{
			Status:  StatusWrongAnswer,
			Comment: res,
		}
	}

	return Verdict{}
}
