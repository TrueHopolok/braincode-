package judge

import (
	"bytes"
	"encoding/gob"
	"errors"
	"strings"

	"github.com/TrueHopolok/braincode-/judge/bf"
)

type bfChecker bf.ByteCode

// NewBFChecker creates a new brainfunk output checker.
//
// # Expected brainfunk API
// Standard input will be test input and test output separated by a 0.
// No output signifies a passes test, any other output will be used as comment.
//
// Not all programs can be tested with this api, for this use [NewLuaChecker], you freak.
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

func (b bfChecker) MarshalBinary() ([]byte, error) { return bf.ByteCode(b).MarshalBinary() }
func (b bfChecker) AppendBinary(buf []byte) ([]byte, error) {
	return bf.ByteCode(b).AppendBinary(buf)
}
func (b *bfChecker) UnmarshalBinary(buf []byte) error {
	return (*bf.ByteCode)(b).UnmarshalBinary(buf)
}
func (b bfChecker) MarshalText() ([]byte, error) { return bf.ByteCode(b).MarshalText() }
func (b bfChecker) AppendText(buf []byte) ([]byte, error) {
	return bf.ByteCode(b).AppendText(buf)
}
func (b *bfChecker) UnmarshalText(buf []byte) error {
	return (*bf.ByteCode)(b).UnmarshalText(buf)
}

type bfSolution struct {
	bc     bf.ByteCode
	steps  int
	memory int
}

// Used for gob marshalling.
type serializedBFSolution struct {
	B bf.ByteCode // Field names are also encoded in GOB, so they are one byte long.
	S int
	M int
}

func (b bfSolution) MarshalBinary() ([]byte, error) { return b.AppendBinary(nil) }
func (b bfSolution) AppendBinary(buf []byte) ([]byte, error) {
	buffer := bytes.NewBuffer(buf)
	err := gob.NewEncoder(buffer).Encode(&serializedBFSolution{
		B: b.bc,
		S: b.steps,
		M: b.memory,
	})
	return buffer.Bytes(), err
}

func (b *bfSolution) UnmarshalBinary(buf []byte) error {
	if b == nil {
		return errors.New("nil receiver")
	}
	var data serializedBFSolution
	if err := gob.NewDecoder(bytes.NewReader(buf)).Decode(&data); err != nil {
		return err
	}
	b.bc = data.B
	b.memory = data.M
	b.steps = data.S
	return nil
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
