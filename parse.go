package main

import "strings"

// blatant copy of github.com/sourcegraph/sourcegraph/internal/own/codeowners/parse.go bits

// isBlank returns true if the current line has no semantically relevant
// content. It can be blank while containing comments or whitespace.
func (p *parsing) isBlank() bool {
	return strings.TrimSpace(p.lineWithoutComments()) == ""
}

const (
	commentStart    = rune('#')
	escapeCharacter = rune('\\')
)

// parsing implements matching and parsing primitives for CODEOWNERS files
// as well as keeps track of internal state as a file is being parsed.
type parsing struct {
	// line is the current line being parsed. CODEOWNERS files are built
	// in such a way that for syntactic purposes, every line can be considered
	// in isolation.
	line string
	// The most recently defined section, or "" if none.
	section string
}

// nextLine advances parsing to focus on the next line.
func (p *parsing) nextLine(line string) {
	p.line = line
}

// lineWithoutComments returns the current line with the commented part
// stripped out.
func (p *parsing) lineWithoutComments() string {
	// A sensible default for index of the first byte where line-comment
	// starts is the line length. When the comment is removed by slicing
	// the string at the end, using the line-length as the index
	// of the first character dropped, yields the original string.
	commentStartIndex := len(p.line)
	var isEscaped bool
	for i, c := range p.line {
		// Unescaped # seen - this is where the comment starts.
		if c == commentStart && !isEscaped {
			commentStartIndex = i
			break
		}
		// Seeing escape character that is not being escaped itself (like \\)
		// means the following character is escaped.
		if c == escapeCharacter && !isEscaped {
			isEscaped = true
			continue
		}
		// Otherwise the next character is definitely not escaped.
		isEscaped = false
	}
	return p.line[:commentStartIndex]
}

func unescape(s string) string {
	var b strings.Builder
	var isEscaped bool
	for _, r := range s {
		if r == escapeCharacter && !isEscaped {
			isEscaped = true
			continue
		}
		b.WriteRune(r)
		isEscaped = false
	}
	return b.String()
}
