package router

import (
	"net/url"
	"strings"
)

func insertToIndex[T any](array []T, index int, value T) []T {
	if len(array) == index { // nil or empty slice or after last element
		return append(array, value)
	}
	array = append(array[:index+1], array[index:]...) // index < len(a)
	array[index] = value
	return array
}

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

func parse(path string) ([]string, url.Values) {
	parts := []string{}
	var builder strings.Builder

	for i, c := range path {
		if c == '?' {
			if builder.Len() > 0 {
				parts = append(parts, builder.String())
			}
			queryParams, err := url.ParseQuery(path[i+1:])
			if err != nil {
				return parts, nil
			}
			return parts, queryParams
		}

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

	return parts, nil
}

type PathPart struct {
	Value      string
	IsVariable bool
}

func parts(path string) []PathPart {
	parts := []PathPart{}
	var builder strings.Builder
	isVariable := false

	for _, c := range path {
		if c == '?' {
			if builder.Len() > 0 {
				parts = append(parts, PathPart{builder.String(), false})
			}
			return parts
		}

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
