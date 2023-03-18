package router

import (
	"net/url"
	"reflect"
	"testing"
)

type pathParseResult struct {
	queryParams url.Values
	parts       []string
}

var pathParseTests = map[string]pathParseResult{
	"/":                               {nil, []string{}},
	"":                                {nil, []string{}},
	"/user":                           {nil, []string{"user"}},
	"/user/profile":                   {nil, []string{"user", "profile"}},
	"/user//profile":                  {nil, []string{"user", "profile"}},
	"/user//Profile":                  {nil, []string{"user", "profile"}},
	"/user//Profile/":                 {nil, []string{"user", "profile"}},
	"/user/profile?id=2":              {url.Values{"id": []string{"2"}}, []string{"user", "profile"}},
	"/product/123/details?sort=likes": {url.Values{"sort": []string{"likes"}}, []string{"product", "123", "details"}},
}

func TestParse(t *testing.T) {
	for path, result := range pathParseTests {
		parts, queryParams := parse(path)
		if !reflect.DeepEqual(result.parts, parts) {
			t.Errorf("Parse error: want %q got %q", result.parts, parts)
		}

		if !reflect.DeepEqual(result.queryParams, queryParams) {
			t.Errorf("Parse error: want %q got %q", result.queryParams, queryParams)
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
