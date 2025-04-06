package judge

import (
	"iter"
	"maps"
)

type listSolution map[string]string

func NewListSolution(answers iter.Seq2[string, string]) OutputChecker {
	return listSolution(maps.Collect(answers))
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
