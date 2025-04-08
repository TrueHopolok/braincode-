package judge

import (
	"bytes"
	"math"
	"strings"
	"sync"

	"github.com/TrueHopolok/braincode-/judge/bf"
)

//go:generate go tool github.com/princjef/gomarkdoc/cmd/gomarkdoc -o documentation.md

// Judge is a handle to a pool of goroutines ready to judge submissions.
//
// Judge must be closed to free up resources.
//
// Judge is safe for concurrent use. Zero value judge is invalid, use [NewJudge] instead.
type Judge struct {
	jobs chan<- job
}

// NewJudge creates a new judge and allocates workers.
// At least one worker will be allocated to prevent a deadlock.
//
// Judge must be closed to free up resources.
func NewJudge(workers int) Judge {
	ch := make(chan job)
	for range max(workers, 1) {
		go worker(ch)
	}
	return Judge{
		jobs: ch,
	}
}

// Frees worker pool of the judge. Never returns an error.
// Signature matches [io.Closer].
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

// Problem is a collection of metadata about a problem. It should be constructed directly.
type Problem struct {
	InputGenerator
	OutputChecker

	Instructions int // Maximum number of active instructions.
	Steps        int // Maximum number of steps of execution.
	Memory       int // Maximum number of allocated bytes during the execution.
}

// CalculateScore is a helper function to calculate score of a given verdict set.
// Returned value is NaN for zero length v and in range [0, 1] in all other cases.
//
// Test group is only counted if all tests in a group pass.
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

// Judge judges a problem against a solution and returns a verdict.
//
// Judge may return a singular verdict [][]Verdict{{value}} if
//   - submition is not valid brainfunk
//   - input generation failed
//   - on any other judge failure (should be unreachable, but who knows)
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
