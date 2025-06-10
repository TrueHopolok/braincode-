// # Syntax
//
// # Structure
//
// MarkLeft documents consists of possible nested blocks. Each block starts with a block directive on a
// new line. It looks like this:
//
//	.name
//	Multi-line data
//	..
//
// Blocks may also be defined inline:
//
//	.name = value of my block. No '..' needed.
//
// Lastly, lines in multi-line blocks may be escaped. If your line starts with a '.' or '!' it has to
// be escaped:
//
//	.myblock
//	!. (escaped  '.')
//	!! (escaped '!')
//	!normal strings can also begin with '!', which will be removed
//	..
//
// If a block must contain other blocks (example: ~C[.quote]), but instead contains raw text, it will
// be interpreted as a paragraph.
//
// # Inline formatting
//
// Paragraphs and section headers support inline formatting. Any style may be nested with any other
// style except itself. Links cannot contain any inline formatting, but may be contained inside other
// styles.
//
//   - bold - ~B[text]
//   - italic - ~I[text]
//   - strike-through - ~S[text]
//   - underline - ~U[text]
//   - code - ~C[text]
//   - math - ~M[text] (LaTex)
//   - link - ~<URL>[text]
//   - escaped ~ - '~~'
//   - escaped ] - '~]'
//
// # Document blocks
//
// '.paragraph' - formatted text. Will be reflown. Cannot contain other blocks.
//
// '.section' - formatted section header. Will be inlined. Cannot contain other blocks.
//
// '.section' - formatted section header. Will be inlined. Cannot contain other blocks.
//
// '.ordered' - ordered list of '.item' blocks. Cannot contain other blocks types.
//
// '.unordered' - unordered list of '.item' blocks. Cannot contain other blocks types.
//
// '.item' - list item. Cannot be used outside lits. May contain any other blocks.
//
// '.quote' - block quote. May contain any other block.
//
// '.code' - code block. Formatting will be preserved. Cannot contain other blocks.
//
// '.example' - code block with sections for input and output. Formatting will be preserved.
// Must contain one '.input' and one '.output' block which act like '.code' blocks.
//
// '.image' - embedded image. Must contain a valid URL.
//
// '.math' - embedded LaTex block.
//
// # Meta blocks
//
// These blocks will not be shown to problem solver directly. Instead they are used to change some
// information about a task. Each of the blocks may be specified only once, with the exception of
// '.task', which can be specified once per locale. Meta blocks may appear anywhere in the document
// tree (except for localization blocks), but preferably should be placed on the top level.
//
// '.task' - title of the solution. Must be plain text, does not support formatting.
//
// '.instructions' - maximum number of instructions submitted solution can have.
//
// '.steps' - maximum number of runtime steps a solution can take.
//
// '.memory' - maximum number of runtime bytes a solution can allocate.
//
// '.[locale]' - mark block for localization. May appear only on the top level (cannot be nested
// inside other blocks). Only locales defined in [KnownLocales] are supported. All blocks outside of
// a localization block belong to a default locale (empty string). Only document blocks and '.task' blocks
// may appear inside localization blocks. Each localization block may be specified multiple times and
// will be equivalent to concatenation of all localization blocks of the same locale.
package ml

//go:generate go tool github.com/princjef/gomarkdoc/cmd/gomarkdoc -o documentation.md

import (
	_ "embed"
	"strings"
	"sync"
)

//go:embed markleft.ml
var docRaw string

var docOnce = sync.OnceValue(func() Document {
	d, err := Parse(strings.NewReader(docRaw))
	if err != nil {
		panic(err)
	}
	return d
})

func Documentation() Document {
	return docOnce()
}
