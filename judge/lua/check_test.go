package lua

import (
	"testing"
)

func TestPrefersCheckerOverSolution(t *testing.T) {
	const source = `
		function checker(input, output)
			if output ~= input .. input then
				return "bad output"
			end
		end

		function solution(input)
			return input .. input
		end

		function test_data()
			return {"1", "2", ""}
		end
	`

	c, err := NewChecker(source)
	if err != nil {
		t.Fatalf("err = %v", err)
	}

	if c.useSolution {
		t.Error("c.useSolution is true, want false")
	}
}
