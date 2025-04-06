package judge

import (
	"bytes"
	"strings"

	"github.com/TrueHopolok/braincode-/judge/bf"
)

type bfChecker bf.ByteCode

func NewBFChecker(source string) (OutputChecker, error) {
	bc, err := bf.Compile(source, -1)
	return bfChecker(bc), err
}

func (b bfChecker) CheckOutput(input string, output string) Verdict {
	buf := new(bytes.Buffer)
	buf.Grow(len(input) + len(output) + 1)
	buf.WriteString(input)
	buf.WriteByte(0)
	buf.WriteString(output)

	out := new(bytes.Buffer)

	s := bf.NewState(bf.ByteCode(b), buf, out, 1e9, 64e6)
	if err := s.Run(); err != nil {
		return Verdict{
			Status:  StatusCheckerFailed,
			Comment: err.Error(),
		}
	}

	if out.Len() > 0 {
		return Verdict{
			Status:  StatusWrongAnswer,
			Comment: out.String(),
		}
	}

	return Verdict{}
}

type bfSolution struct {
	bc     bf.ByteCode
	steps  int
	memory int
}

func NewBFSolution(source string, instructions, steps, memory int) (OutputChecker, error) {
	bc, err := bf.Compile(source, instructions)
	return bfSolution{
		bc:     bc,
		steps:  steps,
		memory: memory,
	}, err
}

func (b bfSolution) CheckOutput(input string, output string) Verdict {
	out := new(bytes.Buffer)
	s := bf.NewState(b.bc, strings.NewReader(input), out, b.steps, b.memory)
	if err := s.Run(); err != nil {
		return Verdict{
			Status:  StatusCheckerFailed,
			Comment: err.Error(),
		}
	}

	if out.String() != output {
		return Verdict{
			Status: StatusWrongAnswer,
		}
	}

	return Verdict{}
}
