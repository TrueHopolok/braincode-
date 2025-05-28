package ml

// Flags for different styles. Due to the syntax rules, it is impossible to parse rich text
// with both [SpanLink]|[SpanCode]
const (
	SpanBold SpanStyle = 1 << iota
	SpanItalic
	SpanCode
	SpanStrike
	SpanUnderline
	SpanLink

	// Number of bits that span [SpanStyle] flags occupy.
	SpanBits = iota
)

type (
	// Document is the root AST node.
	// [Localizations] map must be non nil for any function to work properly.
	Document struct {
		Instructions int
		Steps        int
		Memory       int

		Localizations map[string]*Localizable

		CheckerBF   string
		SolutionBF  string
		GeneratorBF string
		Lua         string
	}

	// Localizable represents all visible localizable content of a document.
	Localizable struct {
		Name   string
		Blocks []Block
	}

	// Block is an closed interface. All block elements implement it.
	Block interface {
		implementsBlock()
	}

	Title RichText

	List struct {
		IsOrdered bool
		Items     []ListItem
	}

	ListItem []Block

	Quote []Block

	Paragraph RichText

	CodeBlock string

	Example struct {
		Input  string
		Output string
	}

	// Image element. String must be a valid URL.
	Image string

	RichText []Span

	SpanStyle byte

	// Span represents a sub string of rich text that has consistent styles.
	//
	// It is impossible for [Style] to be both [SpanLink] and [SpanCode].
	//
	// [URL] is set only if [Style] has bit [SpanLink] set. And vice-versa.
	//
	// [Style] of 0 represents plain text.
	Span struct {
		Style SpanStyle
		Text  string
		URL   string
	}
)

func (Title) implementsBlock()     {}
func (List) implementsBlock()      {}
func (Paragraph) implementsBlock() {}
func (CodeBlock) implementsBlock() {}
func (Example) implementsBlock()   {}
func (Quote) implementsBlock()     {}
func (Image) implementsBlock()     {}
