package ml

import (
	"bufio"
	"cmp"
	"errors"
	"fmt"
	"io"
	"net/url"
	"slices"
	"strconv"
	"strings"
	"unicode/utf8"
)

type blockType int

const (
	// internal, pseudo-blocks
	blockInvalid blockType = iota
	blockDocument
	blockText
	blockLang

	// real blocks
	blockTask
	blockInstructions
	blockSteps
	blockMemory
	blockSection
	blockParagraph
	blockQuote
	blockCode
	blockExample
	blockExampleInput
	blockExampleOutput
	blockImage
	blockChecker
	blockSolution
	blockGenerator
	blockLua
	blockOrdered
	blockUnordered
	blockItem
	blockMath
	blockError
)

var blockToString = map[blockType]string{
	blockTask:         "task",
	blockInstructions: "instructions",
	blockSteps:        "steps",
	blockMemory:       "memory",
	blockSection:      "section",
	blockParagraph:    "paragraph",
	blockQuote:        "quote",
	blockCode:         "code",
	blockExample:      "example",
	blockImage:        "image",
	blockChecker:      "checker",
	blockSolution:     "solution",
	blockGenerator:    "generator",
	blockLua:          "lua",
	blockOrdered:      "ordered",
	blockUnordered:    "unordered",
	blockItem:         "item",
	blockMath:         "math",
	blockError:        "error",
}

var stringToBlock = map[string]blockType{
	"task":         blockTask,
	"instructions": blockInstructions,
	"steps":        blockSteps,
	"memory":       blockMemory,
	"section":      blockSection,
	"paragraph":    blockParagraph,
	"quote":        blockQuote,
	"code":         blockCode,
	"example":      blockExample,
	"input":        blockExampleInput,
	"output":       blockExampleOutput,
	"image":        blockImage,
	"checker":      blockChecker,
	"solution":     blockSolution,
	"generator":    blockGenerator,
	"lua":          blockLua,
	"ordered":      blockOrdered,
	"unordered":    blockUnordered,
	"item":         blockItem,
	"math":         blockMath,
	"error":        blockError,
}

// KnownLocales is a map of all locales. Do not modify concurrently with parsing documents.
//
// Must not be nil.
var KnownLocales = map[string]bool{
	"en": true,
	"ru": true,
}

type parserContext struct {
	buf  *strings.Builder
	path []byte
	errs []error

	Doc           *Document
	CurrentLocale *Localizable
}

func (pctx *parserContext) Buf() *strings.Builder {
	if pctx.buf == nil {
		pctx.buf = new(strings.Builder)
	}
	pctx.buf.Reset()
	return pctx.buf
}

func (pctx *parserContext) PushPath(b *rawBlock) (pop func()) {
	n := len(pctx.path)

	pctx.path = fmt.Appendf(pctx.path, "in block %q (line %d): ", cmp.Or(blockToString[b.kind], b.rawName), b.line)

	isOk := false

	return func() {
		if !isOk {
			return
		}
		pctx.path = pctx.path[:n]
	}
}

func (pctx *parserContext) PushErr(err error) {
	if err == nil {
		return
	}
	pctx.errs = append(pctx.errs, fmt.Errorf("%s%w", pctx.path, err))
}

type blockParser func(pctx *parserContext, b *rawBlock) Block

var blockParsers map[blockType]blockParser

func init() {
	blockParsers = map[blockType]blockParser{
		blockInvalid: func(pctx *parserContext, b *rawBlock) Block {
			pctx.PushErr(errors.New("invalid block"))
			return blockParsers[blockQuote](pctx, b)
		},

		blockLang: func(pctx *parserContext, b *rawBlock) Block {
			pctx.PushErr(errors.New("lang block is only allowed on top level"))
			return blockParsers[blockQuote](pctx, b)
		},

		blockText: func(pctx *parserContext, b *rawBlock) Block {
			return Paragraph(parseRichText(pctx.Buf(), b.data))
		},

		blockTask: func(pctx *parserContext, b *rawBlock) Block {
			pctx.PushErr(stringProperty(pctx.Buf(), &pctx.CurrentLocale.Name, b))
			return nil
		},

		blockInstructions: func(pctx *parserContext, b *rawBlock) Block {
			pctx.PushErr(intProperty(pctx.Buf(), &pctx.Doc.Instructions, b))
			return nil
		},

		blockSteps: func(pctx *parserContext, b *rawBlock) Block {
			pctx.PushErr(intProperty(pctx.Buf(), &pctx.Doc.Steps, b))
			return nil
		},

		blockMemory: func(pctx *parserContext, b *rawBlock) Block {
			pctx.PushErr(intProperty(pctx.Buf(), &pctx.Doc.Memory, b))
			return nil
		},

		blockSection: func(pctx *parserContext, b *rawBlock) Block {
			data, err := textOnly(pctx.Buf(), b)
			if err != nil {
				pctx.PushErr(err)
				return nil
			}
			data = strings.Join(strings.Fields(data), " ")
			if data == "" {
				pctx.PushErr(errors.New("section header cannot be empty"))
				return nil
			}
			return Title(parseRichText(pctx.Buf(), data))
		},

		blockParagraph: func(pctx *parserContext, b *rawBlock) Block {
			data, err := textOnly(pctx.Buf(), b)
			if err != nil {
				pctx.PushErr(err)
				return nil
			}
			data = strings.Join(strings.Fields(data), " ")
			if data == "" {
				pctx.PushErr(errors.New("paragraph cannot be empty"))
				return nil
			}
			return Paragraph(parseRichText(pctx.Buf(), data))
		},

		blockQuote: func(pctx *parserContext, b *rawBlock) Block {
			return Quote(parseBlocks(pctx, b))
		},

		blockCode: func(pctx *parserContext, b *rawBlock) Block {
			data, err := textOnly(pctx.Buf(), b)
			if err != nil {
				pctx.PushErr(err)
				return nil
			}
			return CodeBlock(data)
		},

		blockImage: func(pctx *parserContext, b *rawBlock) Block {
			data, err := textOnly(pctx.Buf(), b)
			if err != nil {
				pctx.PushErr(err)
				return nil
			}

			if _, err := url.Parse(data); err != nil {
				pctx.PushErr(fmt.Errorf("image block must contain a valid URL: %w", err))
				return nil
			}

			return Image(data)
		},

		blockChecker: func(pctx *parserContext, b *rawBlock) Block {
			pctx.PushErr(stringProperty(pctx.Buf(), &pctx.Doc.CheckerBF, b))
			return nil
		},

		blockGenerator: func(pctx *parserContext, b *rawBlock) Block {
			pctx.PushErr(stringProperty(pctx.Buf(), &pctx.Doc.GeneratorBF, b))
			return nil
		},

		blockSolution: func(pctx *parserContext, b *rawBlock) Block {
			pctx.PushErr(stringProperty(pctx.Buf(), &pctx.Doc.SolutionBF, b))
			return nil
		},

		blockLua: func(pctx *parserContext, b *rawBlock) Block {
			pctx.PushErr(stringProperty(pctx.Buf(), &pctx.Doc.Lua, b))
			return nil
		},

		blockOrdered:   parseList(true),
		blockUnordered: parseList(false),

		blockExampleInput: func(pctx *parserContext, b *rawBlock) Block {
			pctx.PushErr(errors.New("input block cannot be used outside of example block"))
			return nil
		},
		blockExampleOutput: func(pctx *parserContext, b *rawBlock) Block {
			pctx.PushErr(errors.New("output block cannot be used outside of example block"))
			return nil
		},
		blockExample: func(pctx *parserContext, b *rawBlock) Block {
			var input, output string
			var inputLine, outputLine int
			for _, c := range b.children {
				pop := pctx.PushPath(c)
				switch c.kind {
				case blockExampleInput:
					if inputLine != 0 {
						pctx.PushErr(fmt.Errorf("duplicate input block (first defined at line %d)",
							inputLine))
					} else {
						s, err := textOnly(pctx.Buf(), c)
						if err != nil {
							pctx.PushErr(err)
						}
						input = s
						inputLine = c.line
					}
				case blockExampleOutput:
					if outputLine != 0 {
						pctx.PushErr(fmt.Errorf("duplicate output block (first defined at line %d)",
							outputLine))
					} else {
						s, err := textOnly(pctx.Buf(), c)
						if err != nil {
							pctx.PushErr(err)
						}
						output = s
						outputLine = c.line
					}
				default:
					pctx.PushErr(errors.New("unsuported block inside example"))
				}
				pop()
			}

			if inputLine == 0 {
				pctx.PushErr(errors.New("missing input block"))
			}
			if outputLine == 0 {
				pctx.PushErr(errors.New("missing output block"))
			}

			return Example{
				Input:  input,
				Output: output,
			}

		},
		blockMath: func(pctx *parserContext, b *rawBlock) Block {
			data, err := textOnly(pctx.Buf(), b)
			if err != nil {
				pctx.PushErr(err)
				return nil
			}
			return Math(data)
		},
	}
}

func parseList(ordered bool) blockParser {
	return func(pctx *parserContext, b *rawBlock) Block {
		var data []ListItem
		for _, c := range b.children {
			pop := pctx.PushPath(c)
			i := ListItem(parseBlocks(pctx, c))
			pop()
			data = append(data, i)
		}
		return List{
			IsOrdered: ordered,
			Items:     data,
		}
	}
}

func stringProperty(buf *strings.Builder, s *string, b *rawBlock) error {
	if *s != "" {
		return errors.New("duplicate definition")
	}
	data, err := textOnly(buf, b)
	if err != nil {
		return err
	}
	if data == "" {
		return errors.New("empty value not allowed")
	}
	*s = data
	return nil
}

func intProperty(buf *strings.Builder, i *int, b *rawBlock) error {
	if *i != 0 {
		return errors.New("duplicate definition")
	}
	data, err := textOnly(buf, b)
	if err != nil {
		return err
	}
	val, err := strconv.Atoi(strings.TrimSpace(data))
	if err != nil {
		return fmt.Errorf("this block must be a valid integer: %w", err)
	}
	*i = val
	return nil
}

func parseBlocks(pctx *parserContext, b *rawBlock) []Block {
	var res []Block
	for _, c := range b.children {
		pop := pctx.PushPath(c)
		parser, found := blockParsers[c.kind]
		if !found {
			panic(fmt.Sprintf("no parser defined for kind %d", c.kind))
		}
		if b := parser(pctx, c); b != nil {
			res = append(res, b)
		}

		pop()
	}
	return res
}

func parseRichText(buf *strings.Builder, s string) RichText {
	const replacement = '\ufffd'

	var res []Span

	if buf == nil {
		buf = new(strings.Builder)
	} else {
		buf.Reset()
	}
	r := []rune(s)

	var state SpanStyle
	var effectStack []SpanStyle
	var lastLink string
	for i := 0; i < len(r); i++ {
		c := r[i]
		if c == 0 || !utf8.ValidRune(c) {
			buf.WriteRune(replacement)
			continue
		}

		if c == ']' {
			if buf.Len() > 0 {
				res = append(res, Span{
					Style: state,
					Text:  buf.String(),
					URL:   lastLink,
				})
				buf.Reset()
			}

			if len(effectStack) == 0 {
				buf.WriteRune(c)
				continue
			}

			last := effectStack[len(effectStack)-1]
			state &^= last
			effectStack = effectStack[:len(effectStack)-1]
			lastLink = ""
			continue
		}

		if c != '~' || i+1 >= len(r) {
			// no tilde / not enough characters
			buf.WriteRune(c)
			continue
		}

		next := r[i+1]
		if r[i+1] == '~' || r[i+1] == ']' {
			// there is some escape
			buf.WriteRune(r[i+1])
			i++
			continue
		}

		if state&(SpanCode|SpanMath|SpanLink) != 0 {
			buf.WriteRune(c)
			continue
		}

		if buf.Len() > 0 {
			res = append(res, Span{
				Style: state,
				Text:  buf.String(),
			})
			buf.Reset()
		}

		var newStyle SpanStyle
		switch next {
		case 'B':
			newStyle = SpanBold
		case 'I':
			newStyle = SpanItalic
		case 'S':
			newStyle = SpanStrike
		case 'U':
			newStyle = SpanUnderline
		case 'C':
			newStyle = SpanCode
		case 'M':
			newStyle = SpanMath
		case '<': // URL
			j := slices.Index(r[i+2:], '>')
			if j == -1 {
				// not terminated link
				buf.WriteRune(c)
				continue
			}

			data := string(r[i+2:][:j])
			if _, err := url.Parse(data); err != nil {
				// invalid link
				buf.WriteRune(c)
				continue
			}

			if j+1 > len(r[i+2:]) || r[i+2:][j+1] != '[' {
				// malformed link text
				buf.WriteRune(c)
				continue
			}
			lastLink = data
			i += j + 3

			state |= SpanLink
			effectStack = append(effectStack, SpanLink)
			continue

		default:
			// bad escape
			buf.WriteRune(c)
			continue
		}

		if i+1 >= len(r) || r[i+2] != '[' {
			// malformed tag
			buf.WriteRune(c)
			continue
		}

		if state&newStyle != 0 {
			// style already set, do nothing
			buf.WriteRune(c)
			continue
		}

		state |= newStyle
		effectStack = append(effectStack, newStyle)

		i += 2
	}

	if buf.Len() > 0 {
		res = append(res, Span{
			Style: state,
			Text:  buf.String(),
		})
		buf.Reset()
	}

	return RichText(res)
}

func textOnly(buf *strings.Builder, b *rawBlock) (string, error) {
	buf.Reset()
	for _, c := range b.children {
		if c.kind != blockText {
			return "", fmt.Errorf("this block cannot contain other blocks, child defined at line %d", c.line)
		}
		buf.WriteString(c.data)
	}
	return buf.String(), nil
}

type rawBlock struct {
	line     int
	kind     blockType
	children []*rawBlock
	data     string // only set for text and locale blocks
	rawName  string
}

type parser struct {
	root         *rawBlock
	state        []*rawBlock
	buf          *strings.Builder
	bufLineStart int
	line         int
	errs         []error
}

func (p *parser) writeLine(line string) {
	p.line++
	s := p.state[len(p.state)-1]

	if len(line) > 0 && line[0] == '!' {
		// escaped line
		if p.buf.Len() == 0 {
			p.bufLineStart = p.line
		}
		p.buf.WriteString(line[1:])
		p.buf.WriteByte('\n')
		return
	}

	if len(line) == 0 {
		return
	}
	if line[0] != '.' {
		// not a directive
		if p.buf.Len() == 0 {
			p.bufLineStart = p.line
		}
		p.buf.WriteString(line)
		p.buf.WriteByte('\n')
		return
	}

	if len(line) > 1 && line[1] == '.' && strings.TrimSpace(line[2:]) == "" {
		// block ended
		if len(p.state) == 1 {
			p.errs = append(p.errs, fmt.Errorf("line %d: extra block closure", p.line))
			return
		}

		if p.buf.Len() > 0 {
			s.children = append(s.children, &rawBlock{
				line: p.line,
				kind: blockText,
				data: p.buf.String(),
			})
			p.buf.Reset()
		}
		p.state = p.state[:len(p.state)-1]
		return
	}

	// a new block directive
	if p.buf.Len() > 0 {
		s.children = append(s.children, &rawBlock{
			kind: blockText,
			data: p.buf.String(),
			line: p.bufLineStart,
		})
		p.buf.Reset()
	}
	i := strings.IndexByte(line, '=')
	if i == -1 {
		// start of a multiline block
		name := strings.TrimSpace(line)[1:]

		if KnownLocales[name] {
			b := &rawBlock{
				line: p.line,
				kind: blockLang,
				data: name,
			}
			s.children = append(s.children, b)
			p.state = append(p.state, b)
			return
		}

		kind, found := stringToBlock[name]
		if !found {
			kind = blockInvalid
		}
		b := &rawBlock{
			kind:    kind,
			line:    p.line,
			rawName: name,
		}
		s.children = append(s.children, b)
		p.state = append(p.state, b)
	} else {
		// inline block
		name := strings.TrimSpace(line[1:i])
		data := strings.TrimSpace(line[i+1:])

		if KnownLocales[name] {
			b := &rawBlock{
				line: p.line,
				kind: blockLang,
				data: name,
				children: []*rawBlock{{
					kind: blockText,
					data: data,
					line: p.line,
				}},
			}
			s.children = append(s.children, b)
			p.state = append(p.state, b)
			return
		}

		kind, found := stringToBlock[name]
		if !found {
			kind = blockInvalid
		}
		s.children = append(s.children, &rawBlock{
			kind: kind,
			children: []*rawBlock{{
				kind: blockText,
				data: data,
				line: p.line,
			}},
			line:    p.line,
			rawName: name,
		})
	}
}

func (p *parser) reset() {
	if p.buf == nil {
		p.buf = new(strings.Builder)
	}
	p.buf.Reset()
	p.state = []*rawBlock{{
		kind: blockDocument,
	}}
	p.bufLineStart = 0
	p.root = p.state[0]
	p.line = 0
}

func (p *parser) finish() (Document, error) {
	if p.buf.Len() > 0 {
		s := p.state[len(p.state)-1]
		s.children = append(s.children, &rawBlock{
			line: p.bufLineStart,
			kind: blockText,
			data: p.buf.String(),
		})
	}
	root := p.root
	buf := p.buf
	pctx := &parserContext{
		buf: buf,
		Doc: &Document{
			Localizations: make(map[string]*Localizable),
		},
	}

	for _, c := range root.children {
		var locale string
		var content []*rawBlock
		if c.kind == blockLang {
			locale = c.data
			content = c.children
		} else {
			locale = ""
			content = []*rawBlock{c}
		}

		if _, found := pctx.Doc.Localizations[locale]; !found {
			pctx.Doc.Localizations[locale] = new(Localizable)
		}
		pctx.CurrentLocale = pctx.Doc.Localizations[locale]

		blocks := parseBlocks(pctx, &rawBlock{
			children: content,
			kind:     blockInvalid,
		})

		pctx.CurrentLocale.Blocks = append(pctx.CurrentLocale.Blocks, blocks...)
	}

	var badLocales []string
	for loc, data := range pctx.Doc.Localizations {
		if data == nil || len(data.Blocks) == 0 {
			badLocales = append(badLocales, loc)
		}
	}
	for _, loc := range badLocales {
		delete(pctx.Doc.Localizations, loc)
	}

	return *pctx.Doc, errors.Join(pctx.errs...)
}

// Parse reads data from r.
//
// Even if error is returned, partial document will still be returned.
func Parse(r io.Reader) (Document, error) {
	p := &parser{}
	p.reset()
	b := bufio.NewScanner(r)
	for b.Scan() {
		p.writeLine(b.Text())
	}
	var ioErr error
	if ioErr = b.Err(); ioErr != nil {
		ioErr = fmt.Errorf("io error at line %d: %w", p.line+1, ioErr)
	}

	doc, err := p.finish()
	err = errors.Join(err, ioErr)

	return doc, err
}

// Format reformats markleft document from r into w.
//
// Note that the document may be loaded into RAM completely.
func Format(r io.Reader, w io.Writer) error {
	d, err := Parse(r)
	if err != nil {
		return err
	}
	return d.WriteSyntax(w)
}
