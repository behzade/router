package router

import (
	"net/http"
	"net/url"
	"reflect"
	"testing"
)

type fakeHandler struct {
	name string
}

func (f *fakeHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
}

type treeTestCase struct {
	handler          http.Handler
	pathWithVariable string
	pathInstance     string
	pathParams       url.Values
	method           string
}

// TODO: better test cases
var insertTests = []treeTestCase{
	{&fakeHandler{"index"}, "", "", url.Values{}, http.MethodGet},
	{&fakeHandler{"profile"}, "/profile", "/profile", url.Values{}, http.MethodGet},
	{&fakeHandler{"profile with id"}, "/profile/{id}", "/profile/2", url.Values{"id": []string{"2"}}, http.MethodGet},
	{&fakeHandler{"register"}, "/register", "/register", url.Values{}, http.MethodPost},
	{&fakeHandler{"user post post"}, "/user/profile/posts/{id}", "/user/profile/posts/3", url.Values{"id": []string{"3"}}, http.MethodPost},
	{&fakeHandler{"user post get"}, "/user/profile/posts/{id}", "/user/profile/posts/3", url.Values{"id": []string{"3"}}, http.MethodGet},
	{&fakeHandler{"user post delete"}, "/user/profile/posts/{id}", "/user/profile/posts/3", url.Values{"id": []string{"3"}}, http.MethodDelete},
	{&fakeHandler{"user post put"}, "/user/profile/posts/{id}", "/user/profile/posts/3", url.Values{"id": []string{"3"}}, http.MethodPut},
	{&fakeHandler{"user post patch"}, "/user/profile/posts/{id}", "/user/profile/posts/3", url.Values{"id": []string{"3"}}, http.MethodPatch},
}

type pathMethod struct {
	path   string
	method string
}

var notFoundTests = []pathMethod{
	{"/asd", http.MethodGet},
	{"/user/profile", http.MethodPost},
	{"/user/posts/asd/22", http.MethodGet},
	{"23/user/posts/asd", http.MethodGet},
	{"/user/posts/", http.MethodGet},
	{"/register/profile", http.MethodGet},
}

var methodNotAllowedTests = []pathMethod{
	{"/", http.MethodPost},
	{"/profile/2", http.MethodDelete},
}

func TestTree(t *testing.T) {
	tree := Tree{map[string]*Tree{}, map[string]*Tree{}, map[string]http.Handler{}}
	var ok bool
	for _, testCase := range insertTests {
		ok = tree.insert(
			parts(testCase.pathWithVariable),
			testCase.method,
			testCase.handler,
		)
		if !ok {
			t.Errorf("Failed to insert route %q", testCase)
		}
	}
	for _, testCase := range insertTests {
		parts, _ := parse(testCase.pathInstance)
		handler, pathParams, statusCode := tree.find(parts, testCase.method)
		if !reflect.DeepEqual(testCase.handler, handler) {
			t.Errorf("Tree find error: expected %q got %q", testCase.handler, handler)
		}

		if !reflect.DeepEqual(testCase.pathParams, pathParams) {
			t.Errorf("Tree find error: expected %q got %q", testCase.pathParams, pathParams)
		}

		if !reflect.DeepEqual(http.StatusOK, statusCode) {
			t.Errorf("Tree find error: expected %q got %q", http.StatusOK, statusCode)
		}
	}

	for _, testCase := range notFoundTests {
		parts, _ := parse(testCase.path)
		handler, pathParams, statusCode := tree.find(parts, testCase.method)
		if statusCode != http.StatusNotFound {
			t.Errorf("Tree find error: expected %q got %q", http.StatusNotFound, statusCode)
		}

		if pathParams != nil {
			t.Errorf("Tree find error: expected %v got %q", nil, pathParams)
		}

		if handler != nil {
			t.Errorf("Tree find error: expected %v got %q", nil, handler)
		}
	}

	for _, testCase := range methodNotAllowedTests {
		parts, _ := parse(testCase.path)
		_, pathParams, statusCode := tree.find(parts, testCase.method)
		if statusCode != http.StatusMethodNotAllowed {
			t.Errorf("Tree find error: expected %q got %q", http.StatusMethodNotAllowed, statusCode)
		}

		if pathParams != nil {
			t.Errorf("Tree find error: expected %v got %q", nil, pathParams)
		}
	}
}
