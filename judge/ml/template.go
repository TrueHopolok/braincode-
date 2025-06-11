package ml

import (
	"context"
	_ "embed"
	"encoding/xml"
	"errors"
	"fmt"
	"html/template"
	"io"
	"maps"
	"net/http"
	"net/url"
	"runtime"
	"slices"
	"strings"
	"sync"
	"time"
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
func (t NestedRichText) IsMath() bool      { return t.Style&SpanMath != 0 }
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
	HTMLClassMathInline

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
	HTMLClassMathBlock

	HTMLClassMapSize
)

// HTMLClassMap contains classes to be emitted into HTML class names.
//
// Quite chonky, should be passed by reference.
type HTMLClassMap [HTMLClassMapSize]string

func (m *HTMLClassMap) TaskTitle() string        { return m[HTMLClassTaskTitle] }
func (m *HTMLClassMap) InfoBlock() string        { return m[HTMLClassInfoBlock] }
func (m *HTMLClassMap) InfoInstructions() string { return m[HTMLClassInfoInstructions] }
func (m *HTMLClassMap) InfoSteps() string        { return m[HTMLClassInfoSteps] }
func (m *HTMLClassMap) InfoMemory() string       { return m[HTMLClassInfoMemory] }
func (m *HTMLClassMap) SpanLink() string         { return m[HTMLClassSpanLink] }
func (m *HTMLClassMap) SpanBold() string         { return m[HTMLClassSpanBold] }
func (m *HTMLClassMap) SpanItalic() string       { return m[HTMLClassSpanItalic] }
func (m *HTMLClassMap) SpanStrike() string       { return m[HTMLClassSpanStrike] }
func (m *HTMLClassMap) SpanCode() string         { return m[HTMLClassSpanCode] }
func (m *HTMLClassMap) SpanUnderline() string    { return m[HTMLClassSpanUnderline] }
func (m *HTMLClassMap) SectionTitle() string     { return m[HTMLClassSectionTitle] }
func (m *HTMLClassMap) Paragraph() string        { return m[HTMLClassParagraph] }
func (m *HTMLClassMap) CodeBlock() string        { return m[HTMLClassCodeBlock] }
func (m *HTMLClassMap) Image() string            { return m[HTMLClassImage] }
func (m *HTMLClassMap) Example() string          { return m[HTMLClassExample] }
func (m *HTMLClassMap) ExampleInput() string     { return m[HTMLClassExampleInput] }
func (m *HTMLClassMap) ExampleOutput() string    { return m[HTMLClassExampleOutput] }
func (m *HTMLClassMap) OrderedList() string      { return m[HTMLClassOrderedList] }
func (m *HTMLClassMap) UnorderedList() string    { return m[HTMLClassUnorderedList] }
func (m *HTMLClassMap) Quote() string            { return m[HTMLClassQuote] }
func (m *HTMLClassMap) ListItem() string         { return m[HTMLClassListItem] }
func (m *HTMLClassMap) MathInline() string       { return m[HTMLClassMathInline] }
func (m *HTMLClassMap) MathBlock() string        { return m[HTMLClassMathBlock] }

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

func AddHTMLTemplate(t *template.Template, name string) {
	t.New(name).Parse(htmlTemplateData)
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
	HTMLClassMathBlock:        "math",
	HTMLClassMathInline:       "inlineMath",
}

const shouldPrerenderMath = runtime.GOARCH != "wasm"

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

	t := TemplatableDocument{
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
	if shouldPrerenderMath {
		t.prerenderMath()
	}
	return t
}

// TemplateContext is a context with helper methods that wraps some other value.
type TemplateContext struct {
	I any

	Debug     bool
	CM        *HTMLClassMap
	MathCache map[string]template.HTML
}

// W returns a [TemplateContext] that wraps i, preserving all other attributes.
func (c TemplateContext) W(i any) TemplateContext {
	c.I = i
	return c
}

// RichTextToNested converts usual, span-based [RichText] into a [NestedRichText].
func (c TemplateContext) RichTextToNested(t RichText) NestedRichText {
	stack := []NestedRichText{{}}
	close := func(style SpanStyle) { // close style
		last := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		prelast := &stack[len(stack)-1]
		prelast.Children = append(prelast.Children, last)
	}
	renderRichRaw(t,
		func(style SpanStyle, url string) { // open style
			stack = append(stack, NestedRichText{
				Style: style,
				URL:   url,
			})
		},
		close,
		func(st SpanStyle, s string) { // print string
			last := &stack[len(stack)-1]

			if st&SpanMath != 0 {
				last.Data += s
			} else {
				last.Children = append(last.Children, NestedRichText{
					Data: s,
				})
			}
		},
	)

	for ; len(stack) > 1; close(0) {
	}

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

func (c TemplateContext) Math() *Math {
	if v, ok := c.I.(Math); ok {
		return &v
	}
	return nil
}

func (c TemplateContext) RenderMathInline(value string) template.HTML {
	return c.renderMath(string(value), c.CM.MathInline())
}

func (c TemplateContext) RenderMathBlock(value Math) template.HTML {
	return c.renderMath(string(value), c.CM.MathBlock())
}

func (c TemplateContext) renderMath(value, class string) template.HTML {
	if res, ok := c.MathCache[value]; ok {
		return res
	}
	return template.HTML(fmt.Sprintf(`<img src="%v" class="%s"/>`, &url.URL{
		Scheme: "https",
		Host:   "latex.codecogs.com",
		Path:   "svg.image",
		// This should be query escape instead, but this API is a bit weird so path escape is used instead.
		RawQuery: url.PathEscape(value),
	}, class))
}

func (c *TemplatableDocument) prerenderMath() {
	if c.MathCache == nil {
		c.MathCache = make(map[string]template.HTML)
	}

	const workers = 4

	type result struct {
		input  string
		output template.HTML
	}

	wg := new(sync.WaitGroup)

	var cMu sync.Mutex

	jobs := make(chan string)
	worker := func() {
		defer wg.Done()

		var res []result

		for source := range jobs {
			svg, err := renderMath(source)
			if err != nil {
				continue
			}
			res = append(res, result{source, svg})
		}

		cMu.Lock()
		for _, v := range res {
			c.MathCache[v.input] = v.output
		}
		cMu.Unlock()
	}

	wg.Add(workers)
	for range workers {
		go worker()
	}

	handled := make(map[string]bool)

	handleRich := func(r RichText) {
		for _, s := range r {
			if s.Style&SpanMath != 0 && !handled[s.Text] {
				s := strings.TrimSpace(string(s.Text))
				handled[s] = true
				jobs <- s
			}
		}
	}

	var walk func(b Block)
	walk = func(b Block) {
		switch v := b.(type) {
		case Math:
			s := strings.TrimSpace(string(v))
			if !handled[s] {
				handled[s] = true
				jobs <- s
			}
		case List:
			for _, item := range v.Items {
				for _, child := range item {
					walk(child)
				}
			}
		case Quote:
			for _, child := range v {
				walk(child)
			}
		case Paragraph:
			handleRich(RichText(v))
		case Title:
			handleRich(RichText(v))
		}
	}

	wg.Add(1)
	for _, b := range c.Blocks {
		walk(b)
	}
	close(jobs)
	wg.Done()

	wg.Wait()
}

func renderMath(data string) (template.HTML, error) {
	u := url.URL{
		Scheme: "https",
		Host:   "latex.codecogs.com",
		Path:   "svg.image",
		// This should be query escape instead, but this API is a bit weird so path escape is used instead.
		RawQuery: url.PathEscape(data),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		cancel()
		return "", err
	}

	req.Header.Add("Accept", "image/svg+xml")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		cancel()
		return "", err
	}
	cancel()

	if resp.StatusCode/100 != 2 {
		return "", errors.New(resp.Status)
	}
	res, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var v struct {
		XMLName xml.Name   `xml:"svg"`
		Attrs   []xml.Attr `xml:",any,attr"`
		Inner   string     `xml:",innerxml"`
	}

	if err := xml.Unmarshal(res, &v); err != nil {
		return "", err
	}

	safe, err := xml.Marshal(&v)
	if err != nil {
		return "", err
	}

	// don't ask...
	return template.HTML(strings.NewReplacer(
		"_xmlns:xlink=", "xmlns:xlink=",
		`xmlns:_xmlns="xmlns"`, "",
	).Replace(string(safe))), nil
}
