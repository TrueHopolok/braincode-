package md

import (
	"bufio"
	"bytes"
	"io"
	"strings"
)

type blockType int

const (
	blockInvalid blockType = iota
	blockDocument
	blockTask
	blockInstructions
	blockSteps
	blockBytes
	blockSection
	blockQuote
	blockCode
	blockExample
	blockImage
	blockChecker
	blockSolution
	blockGenerator
	blockLua
	blockOrdered
	blockUnordered
	blockItem
)

var blockToString = map[blockType]string{
	blockTask:         "task",
	blockInstructions: "instructions",
	blockSteps:        "steps",
	blockBytes:        "bytes",
	blockSection:      "section",
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
}

var stringToBlock = map[string]blockType{
	"task":         blockTask,
	"instructions": blockInstructions,
	"steps":        blockSteps,
	"bytes":        blockBytes,
	"section":      blockSection,
	"quote":        blockQuote,
	"code":         blockCode,
	"example":      blockExample,
	"image":        blockImage,
	"checker":      blockChecker,
	"solution":     blockSolution,
	"generator":    blockGenerator,
	"lua":          blockLua,
	"ordered":      blockOrdered,
	"unordered":    blockUnordered,
	"item":         blockItem,
}

func isWhitespace(s string) bool {
	for i := range s {
		if s[i] != ' ' && s[i] != '\t' {
			return false
		}
	}
	return true
}

func parseBlockHeader(s string) (btype blockType, isInline bool, marker string) {
	btype = blockInvalid
	if len(s) == 0 || s[0] != '.' {
		return
	}
	s = s[1:]
	i := strings.IndexAny(s, "=:")
	if i == -1 {
		return
	}

	name := strings.TrimSpace(s[:i])
	marker = strings.TrimSpace(s[i+1:])

	btype, found := stringToBlock[name]
	if !found {
		return
	}

	isInline = s[i] == '='

	return
}

func isContainer(kind blockType) bool {
	switch kind {
	case blockQuote, blockSection, blockOrdered, blockDocument:
		return true
	default:
		return false
	}
}

func parse(r io.Reader) (Document, error) {
	b := bufio.NewScanner(r)
	var doc Document

	type rawBlock struct {
		kind     blockType
		mark     string
		data     string
		children []*rawBlock
		tag      string
	}

	root := &rawBlock{}
	stack := []*rawBlock{root}
	buf := new(bytes.Buffer)
	for b.Scan() {
		line := b.Text()
		lastBlock := stack[len(stack)-1]

		if lastBlock.tag != "" {
			if !strings.HasPrefix(line, lastBlock.tag) || !isWhitespace(line[len(lastBlock.tag):]) {
				// no tag, keep waiting
				buf.WriteString(line)
				buf.WriteByte('\n')
				continue
			}

			// this block should be closed
			stack = stack[:len(stack)-1]
			lastBlock.data += buf.String()
			buf.Reset()
			line = ""
		}

		canConinue := isContainer(lastBlock.kind)
		t, inline, mark := parseBlockHeader(line)
		if t == blockInvalid {
			// coninue block
			buf.WriteString(line)
			buf.WriteByte('\n')
			continue
		}
		// TODO...

		if len(blocks) == 0 {
			continue
		}

	}

	return doc, b.Err()
}
