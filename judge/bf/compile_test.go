package bf

import "testing"

func Test_stripCommentLoop(t *testing.T) {
	tests := []struct {
		arg     string
		wantRes string
	}{
		{"", ""},
		{"[][][][]+-+-", "+-+-"},
		{"[Hello world!  +++] ", ""},
		{"123 [<>+-[].,] 678  [<>+-[].,]  ", ""},
		{"123 [<>+-].,] 678  [<>+-[].,]", ".,] 678  [<>+-[].,]"},
		{"[[]", "[[]"},
		{"[]]", "[]]"},
		{"asd [[", "asd [["},
		{"[[[][][[[[[]][[]]]]]]]--[+-]", "--[+-]"},
	}
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			if gotRes := stripCommentLoop(tt.arg); gotRes != tt.wantRes {
				t.Errorf("stripCommentLoop() = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}
