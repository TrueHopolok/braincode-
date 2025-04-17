package judge

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"errors"
	"fmt"
	"math"

	"github.com/TrueHopolok/braincode-/judge/lua"
)

type serializedGenerator struct {
	List *listGenerator
	BF   *bfGenerator
	Lua  *luaGenerator
}

func marshalGenerator(enc *gob.Encoder, gen InputGenerator) error {
	var val serializedGenerator
	switch g := gen.(type) {
	case bfGenerator:
		val.BF = &g
	case listGenerator:
		val.List = &g
	case luaGenerator:
		val.Lua = &g
	case *smartGenerator:
		return marshalGenerator(enc, g.InputGenerator)
	default:
		return fmt.Errorf("unexpected generator: %T", g)
	}

	return enc.Encode(&val)
}

func unmarshalGenerator(dec *gob.Decoder) (InputGenerator, error) {
	var val serializedGenerator
	if err := dec.Decode(&val); err != nil {
		return nil, err
	}

	if val.Lua != nil {
		return &smartGenerator{*val.Lua}, nil
	}
	if val.List != nil {
		return *val.List, nil
	}
	if val.BF != nil {
		return &smartGenerator{*val.BF}, nil
	}
	return nil, errors.New("value did not contain any known generator")
}

func MarshalGenerator(g InputGenerator) ([]byte, error) {
	return AppendGenerator(g, nil)
}

func AppendGenerator(g InputGenerator, b []byte) ([]byte, error) {
	buf := bytes.NewBuffer(b)
	err := marshalGenerator(gob.NewEncoder(buf), g)
	return buf.Bytes(), err
}

func UnmarshalGenerator(b []byte) (InputGenerator, error) {
	buf := bytes.NewReader(b)
	res, err := unmarshalGenerator(gob.NewDecoder(buf))
	if err != nil {
		return nil, err
	}
	if buf.Len() > 0 {
		return nil, fmt.Errorf("buffer contains %d trailing junk bytes", buf.Len())
	}
	return res, nil
}

type serializedChecker struct {
	List       *listSolution
	BFSolution *bfSolution
	BFChecker  *bfChecker
	Lua        *lua.Checker
}

func marshalChecker(enc *gob.Encoder, checker OutputChecker) error {
	var val serializedChecker
	switch che := checker.(type) {
	case bfChecker:
		val.BFChecker = &che
	case bfSolution:
		val.BFSolution = &che
	case listSolution:
		val.List = &che
	case *luaChecker:
		val.Lua = &che.Checker
	default:
		return fmt.Errorf("unexpected checker: %T", che)
	}

	return enc.Encode(&val)
}

func unmarshalChecker(dec *gob.Decoder) (OutputChecker, error) {
	var val serializedChecker
	if err := dec.Decode(&val); err != nil {
		return nil, err
	}

	if val.Lua != nil {
		return &luaChecker{*val.Lua}, nil
	}
	if val.BFChecker != nil {
		return *val.BFChecker, nil
	}
	if val.BFSolution != nil {
		return *val.BFSolution, nil
	}
	if val.List != nil {
		return *val.List, nil
	}
	return nil, errors.New("value did not contain any known checker")
}

func MarshalChecker(c OutputChecker) ([]byte, error) {
	return AppendChecker(c, nil)
}

func AppendChecker(c OutputChecker, b []byte) ([]byte, error) {
	buf := bytes.NewBuffer(b)
	err := marshalChecker(gob.NewEncoder(buf), c)
	return buf.Bytes(), err
}

func UnmarshalChecker(b []byte) (OutputChecker, error) {
	buf := bytes.NewReader(b)
	res, err := unmarshalChecker(gob.NewDecoder(buf))
	if err != nil {
		return nil, err
	}
	if buf.Len() > 0 {
		return nil, fmt.Errorf("buffer contains %d trailing junk bytes", buf.Len())
	}
	return res, nil
}

const (
	wireFormatV1 = iota + 1
)

func (p *Problem) MarshalBinary() ([]byte, error) {
	return p.AppendBinary(nil)
}

func (p *Problem) AppendBinary(buf []byte) ([]byte, error) {
	buf = binary.AppendUvarint(buf, wireFormatV1)
	buf = binary.AppendUvarint(buf, uint64(max(p.Instructions, 0)))
	buf = binary.AppendUvarint(buf, uint64(max(p.Memory, 0)))
	buf = binary.AppendUvarint(buf, uint64(max(p.Steps, 0)))

	b := bytes.NewBuffer(buf)
	enc := gob.NewEncoder(b)

	if err := marshalGenerator(enc, p.InputGenerator); err != nil {
		return b.Bytes(), err
	}
	if err := marshalChecker(enc, p.OutputChecker); err != nil {
		return b.Bytes(), err
	}

	return b.Bytes(), nil
}

func (p *Problem) UnmarshalBinary(buf []byte) error {
	r := bytes.NewReader(buf)

	ver, err := binary.ReadUvarint(r)
	if err != nil {
		return err
	}
	if ver != wireFormatV1 {
		return fmt.Errorf("serialized version %v, but parser only recognizes v1", ver)
	}

	instr, err := binary.ReadUvarint(r)
	if err != nil {
		return err
	}
	if instr > math.MaxInt {
		return errors.New("instruction count integer overflow (is this system 32 bit?)")
	}

	memory, err := binary.ReadUvarint(r)
	if err != nil {
		return err
	}
	if memory > math.MaxInt {
		return errors.New("memory limit integer overflow (is this system 32 bit?)")
	}

	steps, err := binary.ReadUvarint(r)
	if err != nil {
		return err
	}
	if steps > math.MaxInt {
		return errors.New("step limit integer overflow (is this system 32 bit?)")
	}

	dec := gob.NewDecoder(r)
	gen, err := unmarshalGenerator(dec)
	if err != nil {
		return err
	}

	che, err := unmarshalChecker(dec)
	if err != nil {
		return err
	}

	if r.Len() > 0 {
		return fmt.Errorf("buffer contains %d trailing junk bytes", r.Len())
	}

	*p = Problem{
		InputGenerator: gen,
		OutputChecker:  che,
		Steps:          int(steps),
		Memory:         int(memory),
		Instructions:   int(instr),
	}
	return nil
}
