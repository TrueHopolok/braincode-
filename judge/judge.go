package judge

import (
	"bytes"
	"fmt"
	"math"
	"strings"
	"sync"

	"github.com/TrueHopolok/braincode-/judge/bf"
)

//go:generate go tool golang.org/x/tools/cmd/stringer -type=Status -trimprefix=Status
type Status uint

const (
	StatusAccept Status = iota
	StatusCompilationFailed
	StatusRuntimeError
	StatusSourceSizeLimit
	StatusTimeLimit
	StatusMemoryLimit
	StatusWrongAnswer
	StatusCheckerFailed
	StatusJudgeFailed
)

type Verdict struct {
	Status  Status
	Comment string
}

func (v Verdict) Error() string {
	if v.Comment != "" {
		return fmt.Sprintf("%v: %s", v.Status, v.Comment)
	}
	return v.Status.String()
}

type InputGenerator interface {
	GenerateInput() ([][]string, error)
}

type OutputChecker interface {
	CheckOutput(input string, output string) Verdict
}

type Judge struct {
	jobs chan<- job
}

func NewJudge(workers int) Judge {
	ch := make(chan job)
	for range workers {
		go worker(ch)
	}
	return Judge{
		jobs: ch,
	}
}

func (j Judge) Close() error {
	close(j.jobs)
	return nil
}

type job struct {
	OutputChecker

	bc     bf.ByteCode
	input  string
	result func(Verdict)

	steps  int
	memory int
}

type Problem struct {
	InputGenerator
	OutputChecker

	Steps        int
	Memory       int
	Instructions int
}

func CalculateScore(v [][]Verdict) float64 {
	var total, good int

outer:
	for _, group := range v {
		total += len(group)

		for _, test := range group {
			if test.Status != StatusAccept {
				continue outer
			}
		}

		good += len(group)
	}

	return float64(good) / float64(total)
}

func (j Judge) Judge(p Problem, submition string) [][]Verdict {
	if p.Memory <= 0 {
		p.Memory = math.MaxInt
	}
	if p.Steps <= 0 {
		p.Steps = math.MaxInt
	}

	bc, err := bf.Compile(submition, p.Instructions)
	if err != nil {
		switch err.(bf.CompilationError).Kind {
		case bf.CompilationInstructionLimit:
			return [][]Verdict{{{
				Status:  StatusSourceSizeLimit,
				Comment: "program exceeds maximum instruction count",
			}}}
		case bf.CompilationUnmatchedParen:
			return [][]Verdict{{{
				Status:  StatusCompilationFailed,
				Comment: err.Error(),
			}}}
		default:
			return [][]Verdict{{{
				Status:  StatusJudgeFailed,
				Comment: err.Error(),
			}}}
		}
	}

	tests, err := p.GenerateInput()
	if err != nil {
		return [][]Verdict{{{
			Status:  StatusCheckerFailed,
			Comment: err.Error(),
		}}}
	}

	res := make([][]Verdict, 0, len(tests))

	wg := new(sync.WaitGroup)
	for _, t := range tests {
		wg.Add(len(t))
		res = append(res, make([]Verdict, len(t)))
	}

	for groupI, group := range tests {
		for testI, inp := range group {
			j.jobs <- job{
				OutputChecker: p.OutputChecker,
				bc:            bc,
				input:         inp,
				result: func(v Verdict) {
					res[groupI][testI] = v
					wg.Done()
				},
				steps:  p.Steps,
				memory: p.Memory,
			}
		}
	}

	wg.Wait()

	return res
}

func worker(ch <-chan job) {
	for j := range ch {
		j.result(judgeTest(j))
	}
}

func judgeTest(j job) Verdict {
	out := new(bytes.Buffer)
	s := bf.NewState(j.bc, strings.NewReader(j.input), out, j.steps, j.memory)

	if err := s.Run(); err != nil {
		return Verdict{
			Status:  StatusRuntimeError,
			Comment: err.Error(),
		}
	}

	return j.CheckOutput(j.input, out.String())
}
