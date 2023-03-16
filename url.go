package router

import "strings"

func SplitPath(path string) ([]string) {
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

