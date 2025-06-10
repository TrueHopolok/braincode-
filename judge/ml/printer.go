package ml

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net/url"
	"slices"
	"strings"
)

type ioError struct{ error }

func makePrinter(w io.Writer) func(format string, args ...any) {
	return func(format string, args ...any) {
		if _, err := fmt.Fprintf(w, format, args...); err != nil {
			panic(ioError{err})
		}
	}
}

func inline(s string) string {
	return strings.Join(strings.Fields(s), " ")
}

func reflow(s string) string {
	const lineLength = 100

	current := 0
	b := new(strings.Builder)
	words := strings.Fields(s)
	if len(words) == 0 {
		return ""
	}
	b.WriteString(words[0])
	current = len(words[0])
	for _, word := range words[1:] {
		if current+len(word)+1 <= lineLength || current == 0 {
			// fits
			b.WriteByte(' ')
			b.WriteString(word)
			current += len(word) + 1
		} else {
			// does not fit
			b.WriteByte('\n')
			b.WriteString(word)
			current = len(word)
		}
	}
	return b.String()
}

func renderRichRaw(
	t RichText,
	opener func(style SpanStyle, url string),
	closer func(style SpanStyle),
	stringer func(style SpanStyle, data string),
) {
	// each span corresponds to index bits.TrailingZeroes8(style)
	// each element denotes length of prefix of this style starting from this index
	lengths := make([][SpanBits]int, len(t)+1)

	for i := len(t) - 1; i >= 0; i-- {
		style := t[i].Style
		for si := range SpanBits {
			s := SpanStyle(1 << si)
			if style&s == 0 {
				continue
			}
			if i >= len(t)-1 || t[i+1].Style&s == 0 {
				lengths[i][si] = 1
			} else {
				lengths[i][si] = lengths[i+1][si] + 1
			}
		}
	}

	var stack []byte
	var prevStyle SpanStyle

	for i, span := range t {
		// close until none non-needed styles are open
		style := prevStyle
		for style != prevStyle&span.Style {
			s := SpanStyle(1 << stack[len(stack)-1])
			stack = stack[:len(stack)-1]
			style &^= s
			closer(s)
		}

		// open the tag in the order of their lengths
		type pair struct {
			si  byte
			len int
		}

		leftovers := span.Style &^ style

		var data [SpanBits]pair
		for si := range byte(SpanBits) {
			s := SpanStyle(1 << si)
			if leftovers&s == 0 {
				data[si].len = 0
				data[si].si = si
			} else {
				data[si].len = lengths[i][si]
				data[si].si = si
			}
		}

		slices.SortFunc(data[:], func(l, r pair) int {
			return r.len - l.len
		})

		for _, p := range data {
			if p.len == 0 {
				break
			}
			s := SpanStyle(1 << p.si)
			stack = append(stack, p.si)
			if s == SpanLink {
				opener(s, span.URL)
			} else {
				opener(s, "")
			}
		}

		stringer(span.Style, span.Text)

		prevStyle = span.Style
	}

	for i := len(stack) - 1; i >= 0; i-- {
		closer(SpanStyle(1 << stack[i]))
	}
}

func renderRich(t RichText) string {
	b := new(strings.Builder)
	escapes := strings.NewReplacer("~", "~~", "]", "~]")
	renderRichRaw(t,
		func(style SpanStyle, url string) {
			switch style {
			case SpanBold:
				b.WriteString("~B[")
			case SpanCode:
				b.WriteString("~C[")
			case SpanItalic:
				b.WriteString("~I[")
			case SpanStrike:
				b.WriteString("~S[")
			case SpanUnderline:
				b.WriteString("~U[")
			case SpanMath:
				b.WriteString("~M[")
			case SpanLink:
				b.WriteString("~<")
				b.WriteString(url)
				b.WriteString(">[")
			default:
				b.WriteString("~invalid[")
			}
		},
		func(style SpanStyle) {
			b.WriteString("]")
		},
		func(style SpanStyle, text string) {
			_, _ = escapes.WriteString(b, text)
		},
	)
	return b.String()
}

func escape(s string) string {
	b := new(strings.Builder)
	for line := range strings.Lines(s) {
		if len(line) > 0 && line[0] == '!' || line[0] == '.' {
			b.WriteByte('!')
		}
		b.WriteString(line)
	}
	return b.String()
}

// WriteSyntax writes formatted markleft source to w.
//
// Error is returned only if [w.Write] fails.
func (d Document) WriteSyntax(w io.Writer) (err error) {
	bw := bufio.NewWriter(w)
	w = bw
	defer func() {
		if res := recover(); res != nil {
			if e, ok := res.(ioError); ok {
				err = errors.Join(err, e.error)
			} else {
				panic(res)
			}
		}
	}()

	printf := makePrinter(w)

	var defaultLocale Localizable

	if d.Localizations != nil && d.Localizations[""] != nil {
		defaultLocale = *d.Localizations[""]
	}
	hasTitle := false
	if defaultLocale.Name != "" {
		hasTitle = true
		printf(".task = %s\n", inline(defaultLocale.Name))
	}
	for locale, data := range d.Localizations {
		if locale != "" && data.Name != "" {
			hasTitle = true
			printf(".%s\n.task = %s\n..\n", locale, inline(data.Name))
		}
	}

	if hasTitle {
		printf("\n")
	}

	printf(".instructions = %d\n", d.Instructions)
	printf(".steps = %d\n", d.Steps)
	printf(".memory = %d\n", d.Memory)

	first := true

	var render func(b Block)
	render = func(b Block) {
		defer func() {
			first = false
		}()

		switch b := b.(type) {
		case CodeBlock:
			printBlock(printf, "code", string(b))

		case Example:
			printf(".example\n")
			printBlock(printf, "input", b.Input)
			printBlock(printf, "output", b.Output)
			printf("..\n")

		case Image:
			u, err := url.Parse(string(b))
			if err != nil {
				u = &url.URL{
					Path: string(b),
				}
			}
			printf(".image = %v\n", u)

		case List:
			if len(b.Items) == 0 {
				if b.IsOrdered {
					printf(".ordered =\n")
				} else {
					printf(".unordered =\n")
				}
				return
			}

			if b.IsOrdered {
				printf(".ordered\n")
			} else {
				printf(".unordered\n")
			}
			first = true
			for _, elem := range b.Items {
				printContainer(printf, render, "item", elem)
				first = false
			}
			printf("..\n")

		case Paragraph:
			printBlock(printf, "paragraph", reflow(renderRich(RichText(b))))

		case Quote:
			first = true
			printContainer(printf, render, "quote", b)

		case Title:
			if !first {
				printf("\n")
			}
			printf(".section = %s\n", inline(renderRich(RichText(b))))

		case Math:
			printBlock(printf, "math", string(b))

		default:
			printf(".error = block of type %T", b)
		}
	}

	printf("\n")

	if len(defaultLocale.Blocks) > 0 {
		for _, b := range defaultLocale.Blocks {
			render(b)
		}
	}

	for locale, data := range d.Localizations {
		if locale == "" || len(data.Blocks) == 0 {
			continue
		}
		printf(".%s\n", locale)
		for _, b := range data.Blocks {
			render(b)
		}
		printf("..\n")
	}

	if d.CheckerBF != "" || d.SolutionBF != "" || d.GeneratorBF != "" || d.Lua != "" {
		printf("\n")
	}

	if d.CheckerBF != "" {
		printBlock(printf, "checker", d.CheckerBF)
	}
	if d.SolutionBF != "" {
		printBlock(printf, "solution", d.SolutionBF)
	}
	if d.GeneratorBF != "" {
		printBlock(printf, "generator", d.GeneratorBF)
	}
	if d.Lua != "" {
		printBlock(printf, "lua", d.Lua)
	}

	return bw.Flush()
}

func printContainer(printf func(string, ...any), render func(Block), name string, blocks []Block) {
	if len(blocks) == 0 {
		printf(".%s =\n", name)
		return
	}
	if len(blocks) == 1 {
		para, ok := blocks[0].(Paragraph)
		if ok {
			printBlock(printf, name, reflow(renderRich(RichText(para))))
			return
		}
	}
	printf(".%s\n", name)
	for _, b := range blocks {
		render(b)
	}
	printf("..\n")
}

func printBlock(printf func(string, ...any), name string, data string) {
	if strings.IndexByte(data, '\n') == -1 {
		printf(".%s = %s\n", name, strings.TrimSpace(data))
		return
	}
	data = strings.TrimSuffix(data, "\n")
	printf(".%s\n%s\n..\n", name, escape(data))
}
