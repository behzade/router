package router

import (
	"reflect"
	"testing"
)

type pathParseResult struct {
	parts []string
}

var pathParseTests = map[string]pathParseResult{
	"/":                    {[]string{}},
	"":                     {[]string{}},
	"/user":                {[]string{"user"}},
	"/user/profile":        {[]string{"user", "profile"}},
	"/user//profile":       {[]string{"user", "profile"}},
	"/user//Profile":       {[]string{"user", "profile"}},
	"/user//Profile/":      {[]string{"user", "profile"}},
	"/product/123/details": {[]string{"product", "123", "details"}},
}

func TestParse(t *testing.T) {
	for path, result := range pathParseTests {
		var parsedPart string
        rest := path
		for _, part := range result.parts {
            parsedPart, rest = parse(rest)
			if part != parsedPart {
				t.Errorf("Parse error: want %q got %q", part, parsedPart)
			}
		}
	}
}

var pathPartsTests = map[string][]PathPart{
	"/":                                    {},
	"/user":                                {PathPart{"user", false}},
	"/user/profile":                        {PathPart{"user", false}, PathPart{"profile", false}},
	"/user/{id}":                           {PathPart{"user", false}, PathPart{"id", true}},
	"v1/user/{id}/profile/posts/{post-id}": {PathPart{"v1", false}, PathPart{"user", false}, PathPart{"id", true}, PathPart{"profile", false}, PathPart{"posts", false}, PathPart{"post-id", true}},
	"{var1}/{var2}/{var3}/{var1}":          {PathPart{"var1", true}, PathPart{"var2", true}, PathPart{"var3", true}, PathPart{"var1", true}},
}

func TestParts(t *testing.T) {
	for path, pathParts := range pathPartsTests {
		parts := parts(path)
		if !reflect.DeepEqual(pathParts, parts) {
			t.Errorf("Parse error: want %v got %v", pathParts, parts)
		}
	}
}
