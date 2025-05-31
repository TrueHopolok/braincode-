// This is a support web assembly module.

//go:build js

package main

import (
	"bytes"
	"fmt"
	"strings"
	"syscall/js"

	"github.com/TrueHopolok/braincode-/judge/ml"
)

type ParseResult struct {
	Source          string
	FormattedSource string
	HTML            string
	Locale          string
	Errors          []string
}

func ParseMarkLeft(source, locale string) ParseResult {
	var errors []string

	doc, err := ml.Parse(strings.NewReader(source))
	if err != nil {
		errs := err.(interface{ Unwrap() []error }).Unwrap()
		errors = make([]string, 0, len(errs))
		for _, err := range errs {
			errors = append(errors, err.Error())
		}
	}

	tpl := ml.HTMLTemplate()
	tdoc := doc.Templatable(locale)

	b := new(bytes.Buffer)
	b.Grow(2048)
	if err := tpl.Execute(b, tdoc); err != nil {
		errors = append(errors, fmt.Sprintf("html template execution failed: %v", err))
	}

	html := b.String()

	b.Reset()
	if err := doc.WriteSyntax(b); err != nil {
		errors = append(errors, fmt.Sprintf("source formatting failed: %v", err))
	}

	return ParseResult{
		Source:          source,
		FormattedSource: b.String(),
		HTML:            html,
		Locale:          locale,
		Errors:          errors,
	}
}

func main() {
	f := js.FuncOf(func(_ js.Value, args []js.Value) any {
		if len(args) == 0 || len(args) > 2 {
			panic(js.Error{
				Value: js.ValueOf(fmt.Sprintf("parseMarkleft: invalid argument count, expected 1 or 2, got %d", len(args))),
			})
		}

		var source string
		var locale string
		source = args[0].String()
		if len(args) > 1 {
			locale = args[1].String()
		}

		pr := ParseMarkLeft(source, locale)
		errs := make([]any, 0, len(pr.Errors))
		for _, err := range pr.Errors {
			fmt.Println(err)
			errs = append(errs, err)
		}

		return map[string]any{
			"source":           pr.Source,
			"formatted_source": pr.FormattedSource,
			"html":             pr.HTML,
			"locale":           pr.Locale,
			"errors":           errs,
		}
	})

	js.Global().Set("parseMarkleft", f)

	select {}
}
