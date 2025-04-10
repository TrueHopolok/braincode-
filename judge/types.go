package judge

import "fmt"

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
