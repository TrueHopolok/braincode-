package bf

import (
	"bufio"
	"errors"
	"io"
)

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

func (s *State) UsedMemory() int {
	return len(s.memory)
}

func (s *State) RemainingSteps() int {
	return s.stepLimit
}

func (s *State) Run() error {
	for ; !s.Finished(); s.Step() {
	}
	return s.Error()
}

func (s *State) Finished() bool {
	return s.instruction >= len(s.bytecode)
}

func (s *State) Error() error {
	return s.err
}

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
