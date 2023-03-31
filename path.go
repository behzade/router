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

// parse path to get usable parts for router and query params
func parse(path string) (string, string) {
	var buf [128]byte
	var i int
	var n int

	for ; i < len(path); i++ {
		if path[i] == '/' && n > 0 {
			return string(buf[:n]), path[i+1:]
		}

		if path[i] >= 'a' && path[i] <= 'z' || path[i] >= '0' && path[i] <= '9' || path[i] == '-' {
			buf[n] = path[i]
			n++
		}

        if path[i] >= 'A' && path[i] <= 'Z' {
			buf[n] = byte(unicode.ToLower(rune(path[i])))
			n++
        }

	}
	return string(buf[:n]), path[i:]
}

type PathPart struct {
	Value      string
	IsVariable bool
}

// split path to variable and constant parts
func parts(path string) []PathPart {
	parts := []PathPart{}
	var builder strings.Builder
	isVariable := false

	for i, c := range path {
		if c == '/' {
			if builder.Len() > 0 {
				parts = append(parts, PathPart{builder.String(), false})
				builder.Reset()
			}
			continue
		}

		if c == '{' && !isVariable {
			isVariable = true
		}

		if c == '}' && isVariable {
			parts = append(parts, PathPart{builder.String(), true})
			isVariable = false
			builder.Reset()
		}

		writeAllowedByte(path[i], &builder)
	}

	if builder.Len() > 0 {
		parts = append(parts, PathPart{builder.String(), false})
	}

	return parts
}
