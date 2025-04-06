package bf

import "fmt"

type CompilationError struct {
	Kind CompilationErrorKind

	Offset int
	Char   rune
}

type CompilationErrorKind int

const (
	CompilationUnmatchedParen CompilationErrorKind = iota
	CompilationInstructionLimit
)

func (e CompilationError) Error() string {
	switch e.Kind {
	case CompilationInstructionLimit:
		return fmt.Sprintf("compilation error: instruction limit, first offending instruction at %d", e.Offset)
	case CompilationUnmatchedParen:
		return fmt.Sprintf("compilation error: bracket '%c' at offset %d is unmatched", e.Char, e.Offset)
	default:
		panic(fmt.Sprintf("unexpected bf.CompilationErrorKind: %#v", e.Kind))
	}
}

type RuntimeError int

const (
	RuntimeErrorHeadUnderflow RuntimeError = iota
	RuntimeErrorMemoryLimit
	RuntimeErrorStepLimit
)

func (e RuntimeError) Error() string {
	switch e {
	case RuntimeErrorHeadUnderflow:
		return "runtime error: head underflow"
	case RuntimeErrorMemoryLimit:
		return "runtime error: memory limit"
	case RuntimeErrorStepLimit:
		return "runtime error: step limit"
	default:
		panic(fmt.Sprintf("unexpected bf.RuntimeError: %#v", e))
	}
}
