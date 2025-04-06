package bf

import (
	"math"
)

type ByteCode struct {
	ops []opjump
}

type Op byte

const (
	OpLeft      Op = '<'
	OpRight     Op = '>'
	OpIncrement Op = '+'
	OpDecrement Op = '-'
	OpInput     Op = ','
	OpOutput    Op = '.'
	OpLoopStart Op = '['
	OpLoopEnd   Op = ']'

	opMax = max(OpLeft, OpRight, OpIncrement, OpDecrement, OpIncrement, OpOutput, OpLoopStart, OpLoopEnd)
)

var validOp = [opMax + 1]bool{
	OpLeft:      true,
	OpRight:     true,
	OpIncrement: true,
	OpDecrement: true,
	OpInput:     true,
	OpOutput:    true,
	OpLoopStart: true,
	OpLoopEnd:   true,
}

type opjump uint32

func makeOpjump(op Op, addr uint32) opjump {
	switch op {
	case OpLoopEnd, OpLoopStart:
		return opjump(addr)<<1 | 1
	default:
		return opjump(op) << 1
	}
}

func (j opjump) Op() Op {
	if j&1 == 0 {
		return Op(j >> 1)
	}
	return OpLoopStart

}

func (j opjump) Addr() uint32 {
	return uint32(j >> 1)

}

func Compile(source string, instructionLimit int) (ByteCode, error) {
	if instructionLimit < 0 {
		instructionLimit = math.MaxInt
	}

	var instructions []opjump
	var openLoops []uint32
	var byteOff []int

	for i, r := range source {
		if r > rune(opMax) || !validOp[r] {
			continue
		}

		if len(instructions) >= instructionLimit {
			return ByteCode{}, CompilationError{
				Kind:   CompilationInstructionLimit,
				Offset: i,
				Char:   r,
			}
		}

		b := Op(r)
		switch b {
		case OpLoopStart:
			openLoops = append(openLoops, uint32(len(instructions)))
			byteOff = append(byteOff, i)
			instructions = append(instructions, makeOpjump(OpLoopStart, 0))

		case OpLoopEnd:
			if len(openLoops) == 0 {
				return ByteCode{}, CompilationError{
					Kind:   CompilationUnmatchedParen,
					Offset: i,
					Char:   rune(OpLoopEnd),
				}
			}
			start := openLoops[len(openLoops)-1]
			openLoops = openLoops[:len(openLoops)-1]
			byteOff = byteOff[:len(byteOff)-1]

			instructions[start] = makeOpjump(OpLoopStart, uint32(len(instructions)))
			instructions = append(instructions, makeOpjump(OpLoopEnd, start))

		default:
			instructions = append(instructions, makeOpjump(b, 0))
		}

	}

	if len(openLoops) > 0 {
		return ByteCode{}, CompilationError{
			Kind:   CompilationUnmatchedParen,
			Offset: byteOff[len(byteOff)-1],
			Char:   rune(OpLoopStart),
		}
	}

	return ByteCode{ops: instructions}, nil
}

func (b ByteCode) String() string {
	buf := make([]byte, 0, len(b.ops))
	for i, op := range b.ops {
		char := op.Op()
		if char != OpLoopStart {
			buf = append(buf, byte(char))
		} else if addr := op.Addr(); int(addr) > i {
			buf = append(buf, byte(OpLoopStart))
		} else {
			buf = append(buf, byte(OpLoopEnd))
		}
	}
	return string(buf)
}

var opIndex = [opMax + 1]byte{
	0:           0,
	OpDecrement: 0,
	OpIncrement: 1,
	OpLeft:      2,
	OpRight:     3,
	OpInput:     4,
	OpOutput:    5,
	OpLoopStart: 6,
	OpLoopEnd:   7,
}

var indexOp = [8]Op{
	0: OpDecrement,
	1: OpIncrement,
	2: OpLeft,
	3: OpRight,
	4: OpInput,
	5: OpOutput,
	6: OpLoopStart,
	7: OpLoopEnd,
}

func (o Op) index() byte {
	return opIndex[o]
}
