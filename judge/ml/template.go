package ml

import (
	_ "embed"
	"html/template"
	"maps"
	"slices"
	"sync"
)

// TemplatableDocument is a friendlier representation of [Document]
// to be used with templates. It has a lot of helper methods.
//
// Can be constructed using [Document.Templatable]
type TemplatableDocument struct {
	Name   string
	Locale string

	Instructions int
	Steps        int
	Memory       int

	Blocks []Block

	CheckerBF   string // Should not be emitted to user.
	SolutionBF  string // Should not be emitted to user.
	GeneratorBF string // Should not be emitted to user.
	Lua         string // Should not be emitted to user.

	TemplateContext
}

// NestedRichText is a nested representation of [RichText]
type NestedRichText struct {
	Style    SpanStyle // Guaranteed to have exactly one bit set.
	Children []NestedRichText
	Data     string // Will be set only if Children is nil.
	URL      string
}

func (t NestedRichText) IsBold() bool      { return t.Style&SpanBold != 0 }
func (t NestedRichText) IsItalic() bool    { return t.Style&SpanItalic != 0 }
func (t NestedRichText) IsStrike() bool    { return t.Style&SpanStrike != 0 }
func (t NestedRichText) IsCode() bool      { return t.Style&SpanCode != 0 }
func (t NestedRichText) IsUnderline() bool { return t.Style&SpanUnderline != 0 }
func (t NestedRichText) IsLink() bool      { return t.Style&SpanLink != 0 }
func (t NestedRichText) IsPlain() bool     { return t.Style == 0 }

const (
	// globals

	HTMLClassTaskTitle = iota
	HTMLClassInfoBlock
	HTMLClassInfoInstructions
	HTMLClassInfoSteps
	HTMLClassInfoMemory

	// inlines

	HTMLClassSpanLink
	HTMLClassSpanBold
	HTMLClassSpanItalic
	HTMLClassSpanStrike
	HTMLClassSpanCode
	HTMLClassSpanUnderline

	// blocks

	HTMLClassSectionTitle
	HTMLClassParagraph
	HTMLClassCodeBlock
	HTMLClassImage
	HTMLClassExample
	HTMLClassExampleInput
	HTMLClassExampleOutput
	HTMLClassOrderedList
	HTMLClassUnorderedList
	HTMLClassListItem
	HTMLClassQuote

	HTMLClassMapSize
)

// HTMLClassMap contains classes to be emitted into HTML class names.
//
// Quite chonky, should be passed by reference.
type HTMLClassMap [HTMLClassMapSize]string

func (m HTMLClassMap) TaskTitle() string        { return m[HTMLClassTaskTitle] }
func (m HTMLClassMap) InfoBlock() string        { return m[HTMLClassInfoBlock] }
func (m HTMLClassMap) InfoInstructions() string { return m[HTMLClassInfoInstructions] }
func (m HTMLClassMap) InfoSteps() string        { return m[HTMLClassInfoSteps] }
func (m HTMLClassMap) InfoMemory() string       { return m[HTMLClassInfoMemory] }
func (m HTMLClassMap) SpanLink() string         { return m[HTMLClassSpanLink] }
func (m HTMLClassMap) SpanBold() string         { return m[HTMLClassSpanBold] }
func (m HTMLClassMap) SpanItalic() string       { return m[HTMLClassSpanItalic] }
func (m HTMLClassMap) SpanStrike() string       { return m[HTMLClassSpanStrike] }
func (m HTMLClassMap) SpanCode() string         { return m[HTMLClassSpanCode] }
func (m HTMLClassMap) SpanUnderline() string    { return m[HTMLClassSpanUnderline] }
func (m HTMLClassMap) SectionTitle() string     { return m[HTMLClassSectionTitle] }
func (m HTMLClassMap) Paragraph() string        { return m[HTMLClassParagraph] }
func (m HTMLClassMap) CodeBlock() string        { return m[HTMLClassCodeBlock] }
func (m HTMLClassMap) Image() string            { return m[HTMLClassImage] }
func (m HTMLClassMap) Example() string          { return m[HTMLClassExample] }
func (m HTMLClassMap) ExampleInput() string     { return m[HTMLClassExampleInput] }
func (m HTMLClassMap) ExampleOutput() string    { return m[HTMLClassExampleOutput] }
func (m HTMLClassMap) OrderedList() string      { return m[HTMLClassOrderedList] }
func (m HTMLClassMap) UnorderedList() string    { return m[HTMLClassUnorderedList] }
func (m HTMLClassMap) Quote() string            { return m[HTMLClassQuote] }
func (m HTMLClassMap) ListItem() string         { return m[HTMLClassListItem] }

//go:embed html.template
var htmlTemplateData string
var htmlParseOnce = sync.OnceValue(func() *template.Template {
	tpl, err := template.New("markleftDoc").Parse(htmlTemplateData)
	if err != nil {
		panic(err)
	}
	return tpl
})

// HTMLTemplate returns a parsed template named "markleftDoc" that expects a [TemplatableDocument] and emits HTML.
func HTMLTemplate() *template.Template {
	tpl, err := htmlParseOnce().Clone()
	if err != nil {
		panic(err)
	}
	return tpl
}

// DefaultClassMap is a the default class map used by parsers if not overridden.
// Should not be modified concurrently with
var DefaultClassMap = HTMLClassMap{
	HTMLClassTaskTitle:        "taskTitle",
	HTMLClassInfoBlock:        "infoBlock",
	HTMLClassInfoInstructions: "infoInstructions",
	HTMLClassInfoSteps:        "infoSteps",
	HTMLClassInfoMemory:       "infoMemory",
	HTMLClassSpanLink:         "spanLink",
	HTMLClassSpanBold:         "spanBold",
	HTMLClassSpanItalic:       "spanItalic",
	HTMLClassSpanStrike:       "spanStrike",
	HTMLClassSpanCode:         "spanCode",
	HTMLClassSpanUnderline:    "spanUnderline",
	HTMLClassSectionTitle:     "sectionTitle",
	HTMLClassParagraph:        "paragraph",
	HTMLClassCodeBlock:        "codeBlock",
	HTMLClassImage:            "image",
	HTMLClassExample:          "example",
	HTMLClassExampleInput:     "exampleInput",
	HTMLClassExampleOutput:    "exampleOutput",
	HTMLClassOrderedList:      "orderedList",
	HTMLClassUnorderedList:    "unorderedList",
	HTMLClassListItem:         "listItem",
	HTMLClassQuote:            "quote",
}

// Make a templatable document from this document.
//
// A locale may not match. In that case, some other available locale will be selected.
// Selected locale can be accessed using [TemplatableDocument.Locale].
func (d Document) Templatable(locale string) TemplatableDocument {
	var loc *Localizable
	if l, ok := d.Localizations[locale]; ok {
		loc = l
	} else if l, ok := d.Localizations[""]; ok {
		locale = ""
		loc = l
	} else {
		keys := slices.Collect(maps.Keys(d.Localizations))
		if len(keys) > 0 {
			key := slices.Min(keys)
			loc = d.Localizations[key]
			locale = key
		} else {
			loc = new(Localizable)
			locale = ""
		}
	}
	return TemplatableDocument{
		Name:         loc.Name,
		Locale:       locale,
		Instructions: d.Instructions,
		Steps:        d.Steps,
		Memory:       d.Memory,
		Blocks:       loc.Blocks,
		TemplateContext: TemplateContext{
			CM: &DefaultClassMap,
		},
	}
}

// TemplateContext is a context with helper methods that wraps some other value.
type TemplateContext struct {
	I any

	Debug bool
	CM    *HTMLClassMap
}

// W returns a [TemplateContext] that wraps i, preserving all other attributes.
func (c TemplateContext) W(i any) TemplateContext {
	c.I = i
	return c
}

// RichTextToNested converts usual, span-based [RichText] into a [NestedRichText].
func (c TemplateContext) RichTextToNested(t RichText) NestedRichText {
	stack := []NestedRichText{{}}
	renderRichRaw(t,
		func(style SpanStyle, url string) {
			stack = append(stack, NestedRichText{
				Style: style,
				URL:   url,
			})
		},
		func(style SpanStyle) {
			stack = stack[:len(stack)-1]
		},
		func(_ SpanStyle, s string) {
			last := &stack[len(stack)-1]
			last.Children = append(last.Children, NestedRichText{
				Data: s,
			})
		},
	)
	return stack[0]
}

func (c TemplateContext) Title() *NestedRichText {
	if v, ok := c.I.(Title); ok {
		res := c.RichTextToNested(RichText(v))
		return &res
	}
	return nil
}

func (c TemplateContext) Paragraph() *NestedRichText {
	if v, ok := c.I.(Paragraph); ok {
		res := c.RichTextToNested(RichText(v))
		return &res
	}
	return nil
}

func (c TemplateContext) CodeBlock() *CodeBlock {
	if v, ok := c.I.(CodeBlock); ok {
		return &v
	}
	return nil
}

func (c TemplateContext) Image() *Image {
	if v, ok := c.I.(Image); ok {
		return &v
	}
	return nil
}

func (c TemplateContext) Example() *Example {
	if v, ok := c.I.(Example); ok {
		return &v
	}
	return nil
}

func (c TemplateContext) List() *List {
	if v, ok := c.I.(List); ok {
		return &v
	}
	return nil
}

func (c TemplateContext) Quote() *Quote {
	if v, ok := c.I.(Quote); ok {
		return &v
	}
	return nil
}
