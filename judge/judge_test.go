package judge_test

import (
	"testing"

	"github.com/TrueHopolok/braincode-/judge"
)

func TestEcho(t *testing.T) {
	const checker = `
	z = string.char(0)

	function checker(input, output)
		if input == output then
			return ""
		else
			return "Expected output to equal input, got: " .. output
		end
	end

	test_data = {
		"hello" .. z,
		"123" .. z,
		"" .. z,
		"Brainfuck!" .. z,
		"a\nb\nc" .. z,
	}`

	const good = `,[.,]>.`
	const bad = `,[,.]`

	runTest(t, checker, good, bad)
}

func TestIncrementAll(t *testing.T) {
	const checker = `
	z = string.char(0)

	function solution(input)
		local expected = input:gsub(".", function(c)
			return string.char((c:byte() + 1) % 256)
		end)
		return expected:sub(1, expected:len()-1) .. z
	end

	test_data = {
		"abc" .. z,
		"123" .. z,
		"" .. z,
		"~}|" .. z,
		"hello world" .. z,
	}`

	const good = `,[+.,]>.`
	const bad = `,[-.,]`

	runTest(t, checker, good, bad)
}

func TestLength(t *testing.T) {
	const checker = `
	function solution(input)
		return string.char(#input - 1)
	end

	test_data = {
		"\0",
		"a\0",
		"ab\0",
		"abc\0",
		"12345\0",
	}
	`

	const good = `,[->+<,]>.`
	const bad = `++++[>++++++++<-],[>+<-]>.`

	runTest(t, checker, good, bad)
}

func runTest(t *testing.T, checker, goodTest, badTest string) {
	t.Helper()

	defer func() {
		t.Helper()
		if err := recover(); err != nil {
			t.Fatal(err)
		}
	}()

	J := judge.NewJudge(1)
	defer J.Close()

	c, err := judge.NewLuaChecker(checker)
	if err != nil {
		t.Fatalf("cannot compile checker: %v", err)
	}

	cc := judge.NewLuaGenerator(checker)

	p := judge.Problem{
		InputGenerator: cc,
		OutputChecker:  c,
		Steps:          1000000,
		Memory:         1000000,
		Instructions:   1000000,
	}

	got := J.Judge(p, goodTest)
	if s := judge.CalculateScore(got); s != 1.0 {
		t.Errorf("got %q (score %v < 1)", got, s)
	}

	got = J.Judge(p, badTest)
	if s := judge.CalculateScore(got); s != 0.0 {
		t.Errorf("got %q (score %v > 0)", got, s)
	}
}
