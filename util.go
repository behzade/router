package router

import "strings"

func insertToIndex[T any](array []T, index int, value T) []T {
	if len(array) == index { // nil or empty slice or after last element
		return append(array, value)
	}
	array = append(array[:index+1], array[index:]...) // index < len(a)
	array[index] = value
	return array
}

func splitPath(path string) []string {
	parts := []string{}
	var builder strings.Builder

	for _, c := range path {
		switch c {
		case '?':
			if builder.Len() > 0 {
				parts = append(parts, builder.String())
			}
			return parts
		case '/':
			if builder.Len() > 0 {
				parts = append(parts, builder.String())
				builder.Reset()
			}
		default:
			builder.WriteRune(c)
		}
	}

	if builder.Len() > 0 {
		parts = append(parts, builder.String())
	}

	return parts
}
