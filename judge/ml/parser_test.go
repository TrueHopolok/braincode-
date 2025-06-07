package ml

import (
	"bytes"
	"io"
	"reflect"
	"strings"
	"testing"

	"github.com/TrueHopolok/braincode-/judge/ml/testhelper"
)

const doc1 = "" +
	`
.task = My task
.section = Hello!
.paragraph
Hello, this is a paragraph.
..
`

func TestParse(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		r       io.Reader
		want    Document
		wantErr bool
	}{
		{"simple", strings.NewReader(doc1), Document{
			Localizations: map[string]*Localizable{
				"": {
					Name: "My task",
					Blocks: []Block{
						Title{Span{Text: "Hello!"}},
						Paragraph{Span{Text: "Hello, this is a paragraph."}},
					},
				}},
		}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := Parse(tt.r)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("Parse() failed: %v", gotErr)
				}
				return
			}

			b := new(bytes.Buffer)
			if err := tt.want.WriteSyntax(b); err != nil {
				t.Error(err)
			}

			if tt.wantErr {
				t.Fatal("Parse() succeeded unexpectedly")
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parse() = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func Test_parseRichText(t *testing.T) {
	type args struct {
		buf *strings.Builder
		s   string
	}
	tests := []struct {
		name string
		args args
		want RichText
	}{
		{"normal", args{nil, `Hello ~B[world]!`}, RichText{
			{Text: "Hello "},
			{Text: "world", Style: SpanBold},
			{Text: "!"},
		}},
		{"escapes", args{nil, `math - ~C[~~M[text~]] (LaTex)`}, RichText{
			{Text: "math - "},
			{Text: "~M[text]", Style: SpanCode},
			{Text: " (LaTex)"},
		}},
		{"math", args{nil, `~M[1 + 1 = 2]`}, RichText{
			{Text: "1 + 1 = 2", Style: SpanMath},
		}},
		{"nested math", args{nil, `~M[1 + ~B[1~] = 2]`}, RichText{
			{Text: "1 + ~B[1] = 2", Style: SpanMath},
		}},
		{"nested code", args{nil, `~C[1 + ~B[1~] = 2]`}, RichText{
			{Text: "1 + ~B[1] = 2", Style: SpanCode},
		}},
		{"nested link", args{nil, `~<http://google.com>[1 + ~B[1~] = 2]`}, RichText{
			{Text: "1 + ~B[1] = 2", Style: SpanLink, URL: "http://google.com"},
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := parseRichText(tt.args.buf, tt.args.s)
			testhelper.DiffValues(t, got, tt.want)
		})
	}
}
