package bf

import "fmt"

// CompilationError is used to report a compilation error.
type CompilationError struct {
	Kind CompilationErrorKind // Kind of error.

	Offset int  // Offending byte offset in the original source code.
	Char   rune // Offending rune. Will always fit into a byte.
}

type CompilationErrorKind int

const (
	CompilationUnmatchedParen   CompilationErrorKind = iota // Loop brackets are not matched.
	CompilationInstructionLimit                             // Compiled program exceeds set limit.
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

// RuntimeError is used to report a failure during execution of a brainfunk program.
type RuntimeError int

const (
	RuntimeErrorHeadUnderflow RuntimeError = iota // Head executed [OpLeft] while pointing to leftmost byte.
	RuntimeErrorMemoryLimit                       // Program allocated limit+1 bytes.
	RuntimeErrorStepLimit                         // Program did not terminate after limit steps.
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
