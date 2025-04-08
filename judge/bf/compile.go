package bf

import (
	"math"
)

// ByteCode contains a compiled brainfunk program.
//
// ByteCode is immutable, so it is safe to use from multiple goroutines.
type ByteCode struct {
	ops []opjump
}

// Op is a brainfunk operations.
type Op byte

const (
	OpLeft      Op = '<' // Move head left. Underflow is a runtime error.
	OpRight     Op = '>' // Move head right. May allocate another byte.
	OpIncrement Op = '+' // Increment value under head. Wraps.
	OpDecrement Op = '-' // Decrement value under head. Wraps.
	OpInput     Op = ',' // Copy byte from input to head.
	OpOutput    Op = '.' // Copy byte from head to output.
	OpLoopStart Op = '[' // Start of loop. Loop acts as a while (head != 0) loop.
	OpLoopEnd   Op = ']' // End of loop.

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

// opjump encode an instruction or a jump offset.
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

// Compile extracts byte code from a source string.
//
// Initial comment loop will be stripped.
//
// Instruction limit may be set to limit length of byte code. Negative limit disables it.
// In that case, [CompilationError] with kind [CompilationInstructionLimit] is never returned.
//
// Returned error's underlying type is [CompilationError].
func Compile(source string, instructionLimit int) (ByteCode, error) {
	if instructionLimit < 0 {
		instructionLimit = math.MaxInt
	}

	was := len(source)
	source = stripCommentLoop(source)
	off := was - len(source)

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
				Offset: i + off,
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
					Offset: i + off,
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
			Offset: byteOff[len(byteOff)-1] + off,
			Char:   rune(OpLoopStart),
		}
	}

	return ByteCode{ops: instructions}, nil
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

func stripCommentLoop(s string) (res string) {
	depth := 0

outer:
	for res = s; len(res) > 0; res = res[1:] {
		b := res[0]
		if res[0] > byte(opMax) || !validOp[b] {
			continue
		}

		switch Op(b) {
		case OpLoopStart:
			depth++
		case OpLoopEnd:
			if depth == 0 {
				return s
			}
			depth--
		default:
			if depth == 0 {
				break outer
			}
		}
	}

	if depth != 0 {
		return s
	}

	return
}
