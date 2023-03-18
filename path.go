package router

import (
	"strings"
	"unicode"
)

func isAllowedChar(c rune) bool {
	if c >= 'a' && c <= 'z' {
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
func parse(path string) []string {
	parts := []string{}
	var builder strings.Builder

	for _, c := range path {
		c = unicode.ToLower(c)

		if c == '/' {
			if builder.Len() > 0 {
				parts = append(parts, builder.String())
				builder.Reset()
			}
			continue
		}

		if isAllowedChar(c) {
			builder.WriteRune(c)
		}
	}

	if builder.Len() > 0 {
		parts = append(parts, builder.String())
	}

	return parts
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

	for _, c := range path {
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

		if isAllowedChar(c) {
			builder.WriteRune(c)
		}
	}

	if builder.Len() > 0 {
		parts = append(parts, PathPart{builder.String(), false})
	}

	return parts
}
