package router

import (
	"strings"
	"unicode"
)

func writeAllowedByte(c byte, b *strings.Builder) {
	if c >= 'a' && c <= 'z' {
		b.WriteByte(c)
	}

	if c >= 'A' && c <= 'Z' {
		b.WriteRune(unicode.ToLower(rune(c)))
	}

	if c >= '0' && c <= '9' {
		b.WriteByte(c)
	}

	if c == '-' {
		b.WriteByte(c)
	}
}

type pathPart struct {
	Value      string
	IsVariable bool
}

type pathParts []pathPart

// split path to variable and constant parts
func parts(path string) []pathPart {
	parts := []pathPart{}
	var builder strings.Builder
	isVariable := false

	for i, c := range path {
		if c == '/' {
			if builder.Len() > 0 {
				parts = append(parts, pathPart{builder.String(), false})
				builder.Reset()
			}
			continue
		}

		if c == '{' && !isVariable {
			isVariable = true
		}

		if c == '}' && isVariable {
			parts = append(parts, pathPart{builder.String(), true})
			isVariable = false
			builder.Reset()
		}

		writeAllowedByte(path[i], &builder)
	}

	if builder.Len() > 0 {
		parts = append(parts, pathPart{builder.String(), false})
	}

	return parts
}

func (p pathParts) dynamicCount() int {
	var i int
	for _, part := range p {
		if part.IsVariable {
			i++
		}
	}
	return i
}
