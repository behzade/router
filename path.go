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
func parse(path string) ([]byte, string) {
	var i int
	var n int

	for ; i < len(path); i++ {
		c := path[i]
		if c == '/' && n > 0 {
			return buf[:n], path[i+1:]
		}

		if c >= 'a' && c <= 'z' || c >= '0' && c <= '9' || c == '-' {
			buf[n] = c
			n++
		} else if c >= 'A' && c <= 'Z' {
			buf[n] = c + 32 // to lower
			n++
		}

	}
	return buf[:n], path[i:]
}

type pathPart struct {
	Value      string
	IsVariable bool
}

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
