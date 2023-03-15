package router

import "strings"

func ParsePath(path string) ([]string, map[string]string) {
	parts := []string{}
	var builder strings.Builder

	for i, c := range path {
		switch c {
		case '?':
			if builder.Len() > 0 {
				parts = append(parts, builder.String())
			}
			return parts, parseQueryString(path[i+1:])
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

	return parts, nil
}

func parseQueryString(queryString string) map[string]string {
	params := map[string]string{}

	isKey := true
	var keyBuilder strings.Builder
	var valueBuilder strings.Builder

	for _, c := range queryString {
		switch c {
		case '=':
			isKey = false
		case '&':
			params[keyBuilder.String()] = valueBuilder.String()
			keyBuilder.Reset()
			valueBuilder.Reset()
			isKey = true
		default:
			if isKey {
				keyBuilder.WriteRune(c)
			} else {
				valueBuilder.WriteRune(c)
			}
		}
	}

	if keyBuilder.Len() > 0 {
		params[keyBuilder.String()] = valueBuilder.String()
	}

	return params
}
