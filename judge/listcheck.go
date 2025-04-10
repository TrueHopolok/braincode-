package judge

import (
	"iter"
	"maps"
)

type listSolution map[string]string

// NewListSolution creates a new table driven test checker. Input must map input to output.
// Output must match exactly.
// Checker will fail for any input not present in answers.
func NewListSolution(answers iter.Seq2[string, string]) OutputChecker {
	return listSolution(maps.Collect(answers))
}

type Pair struct {
	Input  string
	Output string
}

// NewListSolutionSlice is NewListSolution wrapper that does not use iterators.
func NewListSolutionSlice(answers ...Pair) OutputChecker {
	return NewListSolution(pairs(answers).iter())
}

func (l listSolution) CheckOutput(input string, output string) Verdict {
	v, found := l[input]
	if !found {
		return Verdict{
			Status:  StatusCheckerFailed,
			Comment: "checker has no answer",
		}
	}

	if v != output {
		return Verdict{
			Status: StatusWrongAnswer,
		}
	}

	return Verdict{}
}

type pairs []Pair

func (p pairs) iter() iter.Seq2[string, string] {
	return func(yield func(string, string) bool) {
		for _, v := range p {
			if !yield(v.Input, v.Output) {
				break
			}
		}
	}
}
