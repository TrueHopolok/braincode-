// Package bf implements a brainfunk interpreter
//
// Brainfunk needs to be compiled into byte code using [Compile] and then executed using [State].
//
// [Wikipedia]: https://en.wikipedia.org/wiki/Brainfuck
package bf

import (
	"bufio"
	"errors"
	"io"
)

// State stores runtime brainfunk state.
//
// Zero value for state is a terminated program.
//
// State is *not* safe for concurrent use.
type State struct {
	r           io.ByteReader
	w           io.ByteWriter
	stepLimit   int
	memoryLimit int

	bytecode    []opjump
	instruction int

	memory []byte
	head   int

	err error
}

// NewState creates a new [State].
//
// r is used for standard input.
// It may be nil, in that case program will fail with EOF on first read.
// If r implements [io.ByteReader] it is used instead.
//
// w is used for standard output.
// It may be nil, in that program will fail on first write.
// If w implements [io.ByteWriter] it is used instead.
//
// Step and memory limits must be set. Negative limits terminated the program immediately.
func NewState(
	code ByteCode,
	r io.Reader,
	w io.Writer,
	stepLimit int,
	memoryLimit int,
) State {
	if r == nil {
		r = eofReader{}
	}
	if w == nil {
		w = errorWriter{}
	}

	rr, ok := r.(io.ByteReader)
	if !ok {
		rr = bufio.NewReader(r)
	}

	ww, ok := w.(io.ByteWriter)
	if !ok {
		ww = trivialByteWriter{w}
	}

	s := State{
		bytecode:    code.ops,
		r:           rr,
		w:           ww,
		stepLimit:   stepLimit,
		memoryLimit: memoryLimit,
		memory:      []byte{0},
	}
	if memoryLimit < 1 {
		s.finish(RuntimeErrorMemoryLimit)
	}
	return s
}

// UsedMemory returns maximum number of bytes ever used by the program.
//
// Only tape memory is considered, any input or output bytes are not.
// Byte code length is also not counted.
func (s *State) UsedMemory() int {
	return len(s.memory)
}

// RemainingSteps returns maximum number of steps program can take before hitting the step limit.
func (s *State) RemainingSteps() int {
	return s.stepLimit
}

// Run steps repeatedly until program either runs to completion or returns an error.
//
// [Finished] is guaranteed to return true after a call to Run.
//
// See [Step] for more info.
func (s *State) Run() error {
	for ; !s.Finished(); s.Step() {
	}
	return s.Error()
}

// Finished reports whether program completed. All steps will be no-ops after program finished.
func (s *State) Finished() bool {
	return s.instruction >= len(s.bytecode)
}

// Error returns an error if one has happened. If the program did not complete, Error returns a nil error.
func (s *State) Error() error {
	return s.err
}

// Step executes a single brainfunk instruction.
//
// Error is returned in 2 cases:
//   - When a runtime error occurs, in that case underlying type will be [RuntimeError]
//   - If an IO operation fails, in that case error is passes as-is from the reader or writer.
//
// Calling Step on a finished state is a no-op, it will return result of last call to Step.
//
// Any error terminates the brainfunk program.
func (s *State) Step() error {
	if s.Finished() {
		return s.Error()
	}

	if s.stepLimit <= 0 {
		s.finish(RuntimeErrorStepLimit)
		return RuntimeErrorStepLimit
	}
	s.stepLimit--

	switch s.bytecode[s.instruction].Op() {
	case OpDecrement:
		s.memory[s.head]--
	case OpIncrement:
		s.memory[s.head]++
	case OpLeft:
		if s.head == 0 {
			s.finish(RuntimeErrorHeadUnderflow)
			return RuntimeErrorHeadUnderflow
		}
		s.head--

	case OpRight:
		s.head++
		if s.head >= s.memoryLimit {
			s.finish(RuntimeErrorMemoryLimit)
			return RuntimeErrorMemoryLimit
		}
		if len(s.memory) <= s.head {
			s.memory = append(s.memory, 0)
		}

	case OpInput:
		var err error
		s.memory[s.head], err = s.r.ReadByte()
		if err != nil {
			s.finish(err)
			return err
		}

	case OpOutput:
		if err := s.w.WriteByte(s.memory[s.head]); err != nil {
			s.finish(err)
			return err
		}

	case OpLoopStart:
		addr := s.bytecode[s.instruction].Addr()
		if int(addr) > s.instruction {
			// [
			if s.memory[s.head] == 0 {
				// jump over loop
				s.instruction = int(addr)
			}
		} else {
			// ]
			if s.memory[s.head] != 0 {
				// repeat loop
				s.instruction = int(addr)
			}
		}

	default:
		panic("unexpected bf.Op")
	}

	s.instruction++

	return nil
}

func (s *State) finish(err error) {
	s.err = err
	s.instruction = len(s.bytecode)
}

type trivialByteWriter struct{ io.Writer }

func (t trivialByteWriter) WriteByte(c byte) error {
	var buf [1]byte
	buf[0] = c
	_, err := t.Write(buf[:])
	return err
}

type eofReader struct{}

func (eofReader) ReadByte() (byte, error) {
	return 0, io.EOF
}

func (eofReader) Read([]byte) (int, error) {
	return 0, io.EOF
}

type errorWriter struct{}

func (errorWriter) WriteByte(byte) error {
	return errors.New("write to closed output stream")
}

func (errorWriter) Write([]byte) (int, error) {
	return 0, errors.New("write to closed output stream")
}
