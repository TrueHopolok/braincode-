.task = MarkLeft Documentation
.en
.task = MarkLeft Documentation
..
.ru
.task = Документация MarkLeft
..

.instructions = 0
.steps = 0
.memory = 0

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
.item = bold - ~C[~~B[text~]]
.item = italic - ~C[~~I[text~]]
.item = strike-through - ~C[~~S[text~]]
.item = underline - ~C[~~U[text~]]
.item = code - ~C[~~C[text~]]
.item = math - ~C[~~M[text~]] (LaTex)
.item = link - ~C[~~<URL>[text~]]
.item = escaped ~C[~~] - ~C[~~~~]
.item = escaped ~C[~]] - ~C[~~~]]
..

.section = Document blocks
.paragraph = ~C[.paragraph] - formatted text. Will be reflown. Cannot contain other blocks.
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
.paragraph = ~C[.task] - title of the task. Must be plain text, does not support formatting.
.paragraph = ~C[.instructions] - maximum number of instructions submitted solution can have.
.paragraph = ~C[.steps] - maximum number of runtime steps a solution can take.
.paragraph = ~C[.memory] - maximum number of runtime bytes a solution can allocate.
.paragraph
~C[.[locale~]] - mark block for localization. May appear only on the top level (cannot be nested
inside other blocks). Currently only ~C[.ru] and ~C[.en] locales are supported. All blocks outside
of a localization block belong to a ~I[default locale]. Only ~I[Document blocks] and ~C[.task]
blocks may appear inside localization blocks. Each localization block may be specified multiple
times and will be equivalent to concatenation of all localization blocks of the same locale.
..
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
.item = bold - ~C[~~B[text~]]
.item = italic - ~C[~~I[text~]]
.item = strike-through - ~C[~~S[text~]]
.item = underline - ~C[~~U[text~]]
.item = code - ~C[~~C[text~]]
.item = math - ~C[~~M[text~]] (LaTex)
.item = link - ~C[~~<URL>[text~]]
.item = escaped ~C[~~] - ~C[~~~~]
.item = escaped ~C[~]] - ~C[~~~]]
..

.section = Document blocks
.paragraph = ~C[.paragraph] - formatted text. Will be reflown. Cannot contain other blocks.
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
.paragraph = ~C[.task] - title of the task. Must be plain text, does not support formatting.
.paragraph = ~C[.instructions] - maximum number of instructions submitted solution can have.
.paragraph = ~C[.steps] - maximum number of runtime steps a solution can take.
.paragraph = ~C[.memory] - maximum number of runtime bytes a solution can allocate.
.paragraph
~C[.[locale~]] - mark block for localization. May appear only on the top level (cannot be nested
inside other blocks). Currently only ~C[.ru] and ~C[.en] locales are supported. All blocks outside
of a localization block belong to a ~I[default locale]. Only ~I[Document blocks] and ~C[.task]
blocks may appear inside localization blocks. Each localization block may be specified multiple
times and will be equivalent to concatenation of all localization blocks of the same locale.
..
..
.ru

.section = Структура
.paragraph
Документы MarkLeft состоят из потенциально вложенных
блоков. Каждый блок начинается с директивы блока на
отдельной строчке. Это выглядит вот так:
..
.code
!.имя
Многострочные данные.
..
.paragraph = Блоки также могут быть определены на той же строке:
.code
!.имя = значение моего блока. Никаких '..' не надо.
..
.paragraph
И наконец, строки в многострочных блоках могут быть
экранированы. Если ваша строка начинается с символов
'.' или '!' то её нужно экранировать:
..
.code
!.мойблок
!!. (экранированный '.')
!!! (экранированный '!')
!!обычные строки тоже могут быть экранированы, один '!' будет удалён.
!..
..
.paragraph
Если блок должен содержать другие блоки (например
~C[.quote]), но вместо этого содержит текст, текст будет
интерпретирован как параграф.
..

.section = Форматирование текста
.paragraph
Параграфы и зоголовки секций могут содержать
отформатированный текст. Любой стиль может быть
вложен в любой другой, но не в себя же. Ссылка, код и
математический фрагмент не могут содержать другие
стили, но могут сами находится внутри других стилей.
..
.unordered
.item = жирный - ~C[~~B[текст~]]
.item = курсив - ~C[~~I[текст~]]
.item = зачёркнутый - ~C[~~S[текст~]]
.item = подчёркнутый - ~C[~~U[текст~]]
.item = код - ~C[~~C[текст~]]
.item = математика - ~C[~~M[текст~]] (LaTex)
.item = ссылка - ~C[~~<URL>[текст~]]
.item = экранированная ~C[~~] - ~C[~~~~]
.item = экранированная ~C[~]] - ~C[~~~]]
..

.section = Document blocks
.paragraph
~C[.paragraph] - форматированный текст. Пробелы будут
нормализованы. Не может содержать другие блоки.
..
.paragraph
~C[.section] - форматированный заголовок секции. Будет сжат
в одну строчку. Не может содержать другие блоки.
..
.paragraph
~C[.ordered] - упорядоченный список блоков ~C[.item]. Не может
содержать другие типы блоков.
..
.paragraph
~C[.ordered] - неупорядоченный список блоков ~C[.item]. Не может
содержать другие типы блоков.
..
.paragraph
~C[.item] - элемент списка. Не может использоваться вне
списков. Может содержать другие блоки.
..
.paragraph = ~C[.quote] - цитата. Может содержать другие блоки.
.paragraph
~C[.code] - блок кода. Форматирование будет сохранено. Не
может содержать другие блоки.
..
.paragraph
~C[.example] - блок кода с секциями для ввода и вывода.
Форматирование будет сохранено. Обязан содержать
~C[.input] и ~C[.output] блоки, которые ведут себя также как блок
~C[.code].
..
.paragraph = ~C[.image] - встроенная картинка. Должен содержать URL.
.paragraph = ~C[.math] - встроенный блок LaTex математики.

.section = Мета блоки
.paragraph
Эти блоки не будут на прямую показаны участнику. Они
нужны чтобы настраивать те или иные параметры. Каждый
мета блок может использоваться только один раз, за
исключением ~C[.task], который может быть уникальным для
каждой локализации. Мета блоки могут быть где угодно в
документе (за исключением ~I[блоков локализации]), но
обычно они находятся на верхнем уровне.
..
.paragraph
~C[.task] - заголовок задания. Не поддерживает
форматирование.
..
.paragraph
~C[.instructions] - максимальное количество инструкций,
которое может иметь решение.
..
.paragraph
~C[.steps] - максимальное количество шагов выполнения,
которое может иметь решение.
..
.paragraph
~C[.memory] - максимальное количество байтов памяти,
которое может использовать решение.
..
.paragraph
~C[.[locale~]] - маркер локализации. Должен быть на верхнем
уровне (не может быть вложен в другие блоки). На данный
момент поддерживаются только локализации ~C[.ru] и ~C[.en].
Все остальные блоки принадлежат ~I[пустой локализации].
Только ~I[блоки документа] и ~C[.task] могут находится
внутри блока локализации. Блоки локализации могут
быть указаны несколько раз, в таком случае результат
будет равносилен конкатенации каждой локализации по
отдельности.
..
..
