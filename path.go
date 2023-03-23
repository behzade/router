package router

import (
	"strings"
	"unicode"
)

func isAllowedChar(c byte) bool {
	if c >= 'a' && c <= 'z' {
		return true
	}

	if c >= 'A' && c <= 'Z' {
		return true
	}

	if c >= '0' && c <= '9' {
		return true
	}

	if c == '-' {
		return true
	}

	return false
}

// parse path to get usable parts for router and query params
func parse(path string, offset int) (string, int) {
	var builder strings.Builder

	i := offset
	for ; i < len(path); i++ {
		if path[i] == '/' {
			if builder.Len() > 0 {
				return strings.ToLower(builder.String()), i + 1
			}
		}

		if isAllowedChar(path[i]) {
			builder.WriteByte(path[i])
		}
	}
	return strings.ToLower(builder.String()), i
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
		c = unicode.ToLower(c)

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

		if isAllowedChar(path[i]) {
			builder.WriteRune(c)
		}
	}

	if builder.Len() > 0 {
		parts = append(parts, PathPart{builder.String(), false})
	}

	return parts
}
