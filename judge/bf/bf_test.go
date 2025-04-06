package bf_test

import (
	"bytes"
	"reflect"
	"strings"
	"testing"

	"github.com/TrueHopolok/braincode-/judge/bf"
)

func TestHelloWorld(t *testing.T) {
	// thanks wikipedia! (https://en.wikipedia.org/wiki/Brainfuck)
	const source = "++++++++[>++++[>++>+++>+++>+<<<<-]>+>+>->>+[<]<-]>>.>---.+++++++..+++.>>.<-.<.+++.------.--------.>>+.>++."

	bc, err := bf.Compile(source, -1)
	if err != nil {
		t.Fatalf("err = %v", err)
	}

	if s := bc.String(); s != source {
		t.Errorf("have %q want %q", s, source)
	}

	b := new(bytes.Buffer)
	s := bf.NewState(bc, nil, b, 10000, 100)

	if err := s.Run(); err != nil {
		t.Fatalf("err = %v", err)
	}

	const want = "Hello World!\n"
	have := b.String()

	if have != want {
		t.Errorf("have %q want %q", have, want)
	}
}

func TestInput(t *testing.T) {
	const source = `
	[bsort.b -- bubble sort
	(c) 2016 Daniel B. Cristofani
	http://brainfuck.org/]

	>>,[>>,]<<[
	[<<]>>>>[
	<<[>+<<+>-]
	>>[>+<<<<[->]>[<]>>-]
	<<<[[-]>>[>+<-]>>[<<<+>>>-]]
	>>[[<+>-]>>]<
	]<<[>>+<<-]<<
	]>>>>[.>>]

	[This program sorts the bytes of its input by bubble sort.]`

	const input = "abacabadabacaba\x00"

	bc, err := bf.Compile(source, -1)
	if err != nil {
		t.Fatalf("err = %v", err)
	}

	b := new(bytes.Buffer)
	s := bf.NewState(bc, strings.NewReader(input), b, 100000000, 100)

	if err := s.Run(); err != nil {
		t.Fatalf("err = %v", err)
	}

	const want = "aaaaaaaabbbbccd"
	have := b.String()

	if have != want {
		t.Errorf("have %q want %q", have, want)
	}
}

func TestStepLimit(t *testing.T) {
	const source = `
	[e.b -- compute e
	(c) 2016 Daniel B. Cristofani
	http://brainfuck.org/]

	>>>>++>+>++>+>>++<+[  
	[>[>>[>>>>]<<<<[[>>>>+<<<<-]<<<<]>>>>>>]+<]>-
	>>--[+[+++<<<<--]++>>>>--]+[>>>>]<<<<[<<+<+<]<<[
		>>>>>>[[<<<<+>>>>-]>>>>]<<<<<<<<[<<<<]
		>>-[<<+>>-]+<<[->>>>[-[+>>>>-]-<<-[>>>>-]++>>+[-<<<<+]+>>>>]<<<<[<<<<]]
		>[-[<+>-]]+<[->>>>[-[+>>>>-]-<<<-[>>>>-]++>>>+[-<<<<+]+>>>>]<<<<[<<<<]]<<
	]>>>+[>>>>]-[+<<<<--]++[<<<<]>>>+[
		>-[
		>>[--[++>>+>>--]-<[-[-[+++<<<<-]+>>>>-]]++>+[-<<<<+]++>>+>>]
		<<[>[<-<<<]+<]>->>>
		]+>[>>>>]-[+<<<<--]++<[
		[>>>>]<<<<[
			-[+>[<->-]++<[[>-<-]++[<<<<]+>>+>>-]++<<<<-]
			>-[+[<+[<<<<]>]<+>]+<[->->>>[-]]+<<<<
		]
		]>[<<<<]>[
		-[
			-[
			+++++[>++++++++<-]>-.>>>-[<<<----.<]<[<<]>>[-]>->>+[
				[>>>>]+[-[->>>>+>>>>>>>>-[-[+++<<<<[-]]+>>>>-]++[<<<<]]+<<<<]>>>
			]+<+<<
			]>[
			-[
				->[--[++>>>>--]->[-[-[+++<<<<-]+>>>>-]]++<+[-<<<<+]++>>>>]
				<<<<[>[<<<<]+<]>->>
			]<
			]>>>>[--[++>>>>--]-<--[+++>>>>--]+>+[-<<<<+]++>>>>]<<<<<[<<<<]<
		]>[>+<<++<]<
		]>[+>[--[++>>>>--]->--[+++>>>>--]+<+[-<<<<+]++>>>>]<<<[<<<<]]>>
	]>
	]

	This program computes the transcendental number e, in decimal. Because this is
	infinitely long, this program doesn't terminate on its own; you will have to
	kill it. The fact that it doesn't output any linefeeds may also give certain
	implementations trouble, including some of mine.`

	bc, err := bf.Compile(source, -1)
	if err != nil {
		t.Fatalf("err = %v", err)
	}

	b := new(bytes.Buffer)
	s := bf.NewState(bc, nil, b, 100000, 1e8)

	if err := s.Run(); err != bf.RuntimeErrorStepLimit {
		t.Fatalf("err = %v want step limit", err)
	}

	const want = "2.718281"
	have := b.String()

	if have != want {
		t.Errorf("have %q want %q", have, want)
	}
}

func TestMarshal(t *testing.T) {
	// thanks wikipedia! (https://en.wikipedia.org/wiki/Brainfuck)
	const source = "++++++++[>++++[>++>+++>+++>+<<<<-]>+>+>->>+[<]<-]>>.>---.+++++++..+++.>>.<-.<.+++.------.--------.>>+.>++."

	bc, err := bf.Compile(source, -1)
	if err != nil {
		t.Fatalf("err = %v", err)
	}

	bin, err := bc.MarshalBinary()
	if err != nil {
		t.Fatalf("err = %v", err)
	}

	var bc2 bf.ByteCode
	if err := bc2.UnmarshalBinary(bin); err != nil {
		t.Errorf("err = %v", err)
	}

	if !reflect.DeepEqual(bc, bc2) {
		t.Errorf("before %v after %v", bc, bc2)
	}
}
