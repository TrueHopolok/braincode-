package ml

import "testing"

func Test_renderRich(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		t    RichText
		want string
	}{
		{"normal", RichText{
			{SpanBold, "Hello", ""},
			{0, " ", ""},
			{SpanItalic, "world", ""},
			{0, "!", ""},
		}, "~B[Hello] ~I[world]!"},
		{"nested", RichText{
			{SpanBold, "b", ""},
			{SpanItalic, "i", ""},
			{SpanBold | SpanItalic, "bi", ""},
			{0, "n", ""},
		}, "~B[b]~I[i~B[bi]]n"},
		{"link", RichText{
			{SpanBold, "b", ""},
			{SpanItalic, "i", ""},
			{SpanBold | SpanItalic, "bi", ""},
			{SpanBold | SpanItalic | SpanLink, "link", "URL"},
			{0, "n", ""},
		}, "~B[b]~I[i~B[bi~<URL>[link]]]n"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := renderRich(tt.t)
			if got != tt.want {
				t.Errorf("renderRich() = %v, want %v", got, tt.want)
			}
		})
	}
}
