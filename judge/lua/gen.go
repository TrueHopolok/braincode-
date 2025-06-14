package lua

import (
	"context"
	"errors"
	"time"

	lua "github.com/yuin/gopher-lua"
)

const genTimeout = 2 * time.Second

// GetTests extracts tests from lua source.
//
// Tests can be provided in of 2 ways.
//
// # Static tests
//
//	# singular test
//	test_data = "abc"
//
//	# multiple tests without groups
//	test_data = {"1", "2", "3"}
//
//	# multiple tests with subgroup
//	test_data = {{"1", "2"}, {"3", "4"}, nil, "5"}
//
//	# same as this
//	test_data = {{"1", "2"}, {"3", "4"}, {"5"}}
//
// Only strings and nils are allowed. Nils are silently dropped.
//
// # Dynamic tests
//
// A function called test_data may be provided. It will be called with zero arguments and must return
// a string | (nil | string | (nil | string)[])[]
func GetTests(source string) ([][]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), genTimeout)
	defer cancel()

	l := newLuaState()

	defer l.Close()
	l.SetContext(ctx)

	lua.OpenBase(l.LState)

	if err := l.DoString(source); err != nil {
		return nil, err
	}

	val := l.GetGlobal("test_data")

	if val.Type() == lua.LTFunction {
		if err := l.CallByParam(lua.P{
			Fn:      val,
			NRet:    1,
			Protect: true,
		}); err != nil {
			return nil, err
		}

		val = l.Get(-1)
		l.Pop(1)
	}

	return parseTests(val)
}

func parseTests(v lua.LValue) ([][]string, error) {
	switch v := v.(type) {
	case lua.LString:
		return [][]string{{string(v)}}, nil
	case *lua.LTable:
		var res [][]string
		for i := 1; i <= min(v.MaxN(), 100_000); i++ {
			vv := v.RawGetInt(i)
			switch vv := vv.(type) {
			case *lua.LNilType:
				continue

			case lua.LString:
				res = append(res, []string{string(vv)})

			case *lua.LTable:
				var group []string
				for j := 1; j <= min(v.MaxN(), 100_000); j++ {
					vvv := vv.RawGetInt(j)
					switch vvv := vvv.(type) {
					case *lua.LNilType:
						continue
					case lua.LString:
						group = append(group, string(vvv))
					default:
						return nil, errors.New("nested table contains non string test")
					}
				}
				if len(group) > 0 {
					res = append(res, group)
				}

			default:
				return nil, errors.New("table elements' types are invalid")
			}
		}

		if len(res) == 0 {
			return nil, errors.New("empty test data")
		}

		return res, nil

	default:
		return nil, errors.New("invalid tests_data type")
	}

}
