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

var pathPartsTests = map[string][]pathPart{
	"/":                                    {},
	"/user":                                {pathPart{"user", false}},
	"/user/profile":                        {pathPart{"user", false}, pathPart{"profile", false}},
	"/user/{id}":                           {pathPart{"user", false}, pathPart{"id", true}},
	"v1/user/{id}/profile/posts/{post-id}": {pathPart{"v1", false}, pathPart{"user", false}, pathPart{"id", true}, pathPart{"profile", false}, pathPart{"posts", false}, pathPart{"post-id", true}},
	"{var1}/{var2}/{var3}/{var1}":          {pathPart{"var1", true}, pathPart{"var2", true}, pathPart{"var3", true}, pathPart{"var1", true}},
}

func TestParts(t *testing.T) {
	for path, pathParts := range pathPartsTests {
		parts := parts(path)
		if !reflect.DeepEqual(pathParts, parts) {
			t.Errorf("Parse error: want %v got %v", pathParts, parts)
		}
	}
}
