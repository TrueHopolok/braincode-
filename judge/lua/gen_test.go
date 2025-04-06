package lua_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/TrueHopolok/braincode-/judge/lua"
)

func TestGetTests(t *testing.T) {
	tests := []struct {
		name    string
		source  string
		want    [][]string
		wantErr bool
	}{
		{
			"constant string",
			`test_data = "this is my test"`,
			[][]string{{"this is my test"}},
			false,
		},
		{
			"constant flat table",
			`test_data = {"foo", "bar", nil, {"baz1", "baz2"}, {}}`,
			[][]string{{"foo"}, {"bar"}, {"baz1", "baz2"}},
			false,
		},
		{
			"constant nested table",
			`test_data = {"foo", "bar", nil, {"baz1", "baz2"}, {}}`,
			[][]string{{"foo"}, {"bar"}, {"baz1", "baz2"}},
			false,
		},
		{
			"function",
			`function test_data() 
				local result = {}
				for i = 1, 10 do
					table.insert(result, tostring(i))
				end
				return result
			end`,
			[][]string{{"1"}, {"2"}, {"3"}, {"4"}, {"5"}, {"6"}, {"7"}, {"8"}, {"9"}, {"10"}},
			false,
		},
		{
			"empty",
			`test_data = {{}, {}, nil, {}}`,
			nil,
			true,
		},
		{
			"empty",
			`test_data = {}`,
			nil,
			true,
		},
		{
			"bad type",
			`test_data = 123`,
			nil,
			true,
		},
		{
			"infinite loop",
			`function test_data() 
				while true do
					print("lol")
				end
			end`,
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := lua.GetTests(tt.source)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetTests() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetTests() = %v, want %v", got, tt.want)
			}
			fmt.Println(err)
		})
	}
}
