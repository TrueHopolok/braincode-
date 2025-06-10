package ml

import (
	"bytes"
	_ "embed"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/TrueHopolok/braincode-/judge/ml/testhelper"
)

func TestFull(t *testing.T) {
	const data = `
.task = My task!

.memory = 20
.instructions = 10
.steps = 30

.paragraph
I though long and hard and bla bla bla
..

.quote
.paragraph
Hello!
..
.paragraph
Hello 2!
..
.paragraph
Hello 3!
..
..

.paragraph
Here is an example:
..
.example
.input
bla bla bla
..
.output = hello world!
..

.solution = [BRAINFUNK SOLUTION TO THE PROBLEM]
.generator
[OTHER BRAINFUNK PROGRAM]
..
`
	doc, err := Parse(strings.NewReader(data))
	if err != nil {
		t.Fatalf("cannot parse document, err = %v", err)
	}

	want := Document{
		Instructions: 10,
		Steps:        30,
		Memory:       20,
		Localizations: map[string]*Localizable{
			"": {
				Name: "My task!",
				Blocks: []Block{
					Paragraph{{Text: "I though long and hard and bla bla bla"}},
					Quote{
						Paragraph{{Text: "Hello!"}},
						Paragraph{{Text: "Hello 2!"}},
						Paragraph{{Text: "Hello 3!"}},
					},
					Paragraph{{Text: "Here is an example:"}},
					Example{
						Input:  "bla bla bla\n",
						Output: "hello world!",
					},
				},
			},
		},
		SolutionBF:  "[BRAINFUNK SOLUTION TO THE PROBLEM]",
		GeneratorBF: "[OTHER BRAINFUNK PROGRAM]\n",
	}

	testhelper.DiffValues(t, doc, want)

	fmt.Println()
	fmt.Println()
	if err := doc.WriteSyntax(os.Stdout); err != nil {
		t.Error(err)
	}
	fmt.Println()
	fmt.Println()

	tpl := HTMLTemplate()
	td := doc.Templatable("")
	td.Debug = true

	if err := tpl.Execute(os.Stdout, td); err != nil {
		t.Error(err)
	}
}

func TestDocumentation(t *testing.T) {
	doc := Documentation()

	b := new(bytes.Buffer)
	if err := doc.WriteSyntax(b); err != nil {
		t.Fatalf("cannot write syntax: %v", err)
	}

	{
		reparsed, err := Parse(bytes.NewReader(b.Bytes()))
		if err != nil {
			t.Errorf("cannot parse generated source: %v", err)
		} else {
			testhelper.DiffValues(t, doc, reparsed)
		}
	}

	td := doc.Templatable("")
	tpl := HTMLTemplate()

	bbb := new(bytes.Buffer)
	if err := tpl.Execute(bbb, td); err != nil {
		t.Error(err)
	}

	// os.WriteFile("output.html", bbb.Bytes(), 0666)
}
