package md

/*


.task = NAME BLA BLA BLA
.instructions = 120
.steps = 10
.bytes = 234

.section = Hello!

.paragraph: >>
 Paragraph
 as
 para
 para

.ordered:MARKER
.item:
.item:
.item:
MARKER

quote:
end:

code:MARKER
MARKER
~~~X
~~~X

<<<X
Input
~~~X
AAaaa
>>>X

paragraph:  B{bold} I{italic} C{code}  S{strikethrough}

CHECKR

========
group:
test:X
test (%):

end:X
===============
checker:

===============
solution:

===============
generator:

========
lua:

*/

const (
	SpanBold SpanStyle = 1 << iota
	SpanItalic
	SpanCode
	SpanStrike
)

type (
	Document struct {
		Name string

		Instructions int
		Steps        int
		Memory       int

		Blocks []Block

		Tests       [][]string
		CheckerBF   string
		SolutionBF  string
		GeneratorBF string
		Lua         string
	}

	Block interface{}

	Title RichText

	List struct {
		IsOrdered bool
		Items     []Block
	}

	ListItem []Block

	Paragraph RichText

	CodeBlock string

	Example struct {
		Input  string
		Output string
	}

	RichText []Span

	SpanStyle byte

	Span struct {
		Style SpanStyle
		Text  string
	}
)
