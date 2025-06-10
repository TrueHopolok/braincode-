package judge_test

import (
	"strings"
	"testing"

	"github.com/TrueHopolok/braincode-/judge"
	"github.com/TrueHopolok/braincode-/judge/ml"
)

func TestDocumentAPlusB(t *testing.T) {
	const data = `
.task = A + B

.steps = 10000
.instructions = 100
.memory = 200

.en
.paragraph = Input: 2 unsigned bytes: ~C[a] and ~C[b]. Output: ~C[(a + b) mod 256].
..

.ru
.paragraph = Ввод: 2 безнаковых байта: ~C[a] и ~C[b]. Вывод: ~C[(a + b) mod 256].
..

.lua
function solution(input)
	local a = string.byte(input, 1)
	local b = string.byte(input, 2)
	return string.char((a + b) % 256)
end

test_data = {
	{
		string.char(0) .. string.char(0),
		string.char(0) .. string.char(1),
		string.char(0) .. string.char(2),
		string.char(0) .. string.char(3),
		string.char(1) .. string.char(0),
		string.char(1) .. string.char(1),
		string.char(1) .. string.char(2),
		string.char(1) .. string.char(3),
		string.char(2) .. string.char(0),
		string.char(2) .. string.char(1),
		string.char(2) .. string.char(2),
		string.char(2) .. string.char(3),
	},
	{
		string.char(5) .. string.char(6),
		string.char(123) .. string.char(22),
		string.char(55) .. string.char(11),
		string.char(1) .. string.char(3),
		string.char(2) .. string.char(0),
		string.char(2) .. string.char(1),
		string.char(2) .. string.char(2),
	},
	{
		string.char(0) .. string.char(1),
		string.char(255) .. string.char(255),
		string.char(1) .. string.char(3),
		string.char(2) .. string.char(0),
		string.char(2) .. string.char(1),
		string.char(2) .. string.char(2),
	}
}
..
`

	const (
		submitionGood = `
,
>,
<
[->+<]
>.
`
		submitionBad = `
,
>,
<
[->+<]
>+.
`
	)

	doc, err := ml.Parse(strings.NewReader(data))
	if err != nil {
		t.Fatal(err)
	}

	p, err := judge.NewProblem(doc)
	if err != nil {
		t.Fatal(err)
	}

	j := judge.NewJudge(1)
	defer j.Close()

	vGood := j.Judge(p, submitionGood)

	// check that all results are good
	for gi, group := range vGood {
		for ti, v := range group {
			if v.Status != judge.StatusAccept {
				t.Errorf("group %d test %d: %v (%v)", gi, ti, v.Status, v.Comment)
			}
		}
	}

	vBad := j.Judge(p, submitionBad)
	t.Logf("Bad solutions: %v", vBad)

	// check that all results are bad
	for gi, group := range vBad {
		for ti, v := range group {
			if v.Status != judge.StatusWrongAnswer {
				t.Errorf("group %d test %d: %v (%v) - want %v", gi, ti, v.Status, v.Comment, judge.StatusWrongAnswer)
			}
		}
	}

}
