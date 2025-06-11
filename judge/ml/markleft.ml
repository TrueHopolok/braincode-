.en
.task = MarkLeft
..

.instructions = 0
.steps = 0
.memory = 0

.en
.section = Structure
.paragraph
MarkLeft documents consists of possible nested blocks. Each block starts with a block directive on a
new line. It looks like this:
..
.code
!.name
Multi-line data
..
.paragraph = Blocks may also be defined inline:
.code
!.name = value of my block. No '..' needed.
..
.paragraph
Lastly, lines in multi-line blocks may be escaped. If your line starts with a '.' or '!' it has to
be escaped:
..
.code
!.myblock
!!. (escaped '.')
!!! (escaped '!')
!!normal strings can also begin with '!', which will be removed
!..
..
.paragraph
If a block must contain other blocks (example: ~C[.quote]), but instead contains raw text, it will
be interpreted as a paragraph.
..

.section = Inline formatting
.paragraph
Paragraphs and section headers support inline formatting. Any style may be nested with any other
style except itself. Links cannot contain any inline formatting, but may be contained inside other
styles.
..
.unordered
.item = bold - ~C[~B[text~]]
.item = italic - ~C[~I[text~]]
.item = strike-through - ~C[~S[text~]]
.item = underline - ~C[~U[text~]]
.item = code - ~C[~C[text~]] 
.item = math - ~C[~M[text~]] (LaTex)
.item = link - ~C[~<URL>[text~]]
.item = escaped ~C[~~] - ~C[~~~~]
.item = escaped ~C[~]] - ~C[~~~]]
..

.section = Document blocks
.paragraph = ~C[.paragraph] - formatted text. Will be reflown. Cannot contain other blocks.
.paragraph = ~C[.section] - formatted section header. Will be inlined. Cannot contain other blocks.
.paragraph = ~C[.section] - formatted section header. Will be inlined. Cannot contain other blocks.
.paragraph = ~C[.ordered] - ordered list of ~C[.item] blocks. Cannot contain other blocks types.
.paragraph = ~C[.unordered] - unordered list of ~C[.item] blocks. Cannot contain other blocks types.
.paragraph = ~C[.item] - list item. Cannot be used outside lits. May contain any other blocks.
.paragraph = ~C[.quote] - block quote. May contain any other block.
.paragraph = ~C[.code] - code block. Formatting will be preserved. Cannot contain other blocks.
.paragraph
~C[.example] - code block with sections for input and output. Formatting will be preserved. Must
contain ~C[.input] and ~C[.output] blocks which act like ~C[.code] blocks.
..
.paragraph = ~C[.image] - embedded image. Must contain a valid URL.
.paragraph = ~C[.math] - embedded LaTex block.

.section = Meta blocks
.paragraph
These blocks will not be shown to problem solver directly. Instead they are used to change some
information about a task. Each of the blocks may be specified only once, with the exception of
~C[.task], which can be specified once per locale. Meta blocks may appear anywhere in the document
tree (except for ~I[localization blocks]), but preferably should be placed on the top level.
..
.paragraph = ~C[.task] - title of the solution. Must be plain text, does not support formatting.
.paragraph = ~C[.instructions] - maximum number of instructions submitted solution can have.
.paragraph = ~C[.steps] - maximum number of runtime steps a solution can take.
.paragraph = ~C[.memory] - maximum number of runtime bytes a solution can allocate.
.paragraph
~C[.[locale~]] - mark block for localization. May appear only on the top level (cannot be nested
inside other blocks). Currently only ~C[.ru] and ~C[.en] locales are supported. All blocks outside 
of a localization block belong to a -I[default locale~]. Only ~I[Document blocks] and ~C[.task] 
blocks may appear inside localization blocks. Each localization block may be specified multiple 
times and will be equivalent to concatenation of all localization blocks of the same locale.
..
