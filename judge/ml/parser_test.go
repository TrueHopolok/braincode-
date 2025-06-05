package ml_test

import (
	"bytes"
	"io"
	"reflect"
	"strings"
	"testing"

	"github.com/TrueHopolok/braincode-/judge/ml"
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
		want    ml.Document
		wantErr bool
	}{
		{"simple", strings.NewReader(doc1), ml.Document{
			Localizations: map[string]*ml.Localizable{
				"": {
					Name: "My task",
					Blocks: []ml.Block{
						ml.Title{ml.Span{Text: "Hello!"}},
						ml.Paragraph{ml.Span{Text: "Hello, this is a paragraph."}},
					},
				}},
		}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := ml.Parse(tt.r)
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
			// fmt.Println(b.String())

			if tt.wantErr {
				t.Fatal("Parse() succeeded unexpectedly")
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parse() = %+v, want %+v", got, tt.want)
			}
		})
	}
}
