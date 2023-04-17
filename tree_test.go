package router

import (
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"testing"
)

func fakeHandler(http.ResponseWriter, *http.Request, url.Values) {
}

type treeTestCase struct {
	handler          Handler
	pathWithVariable string
	pathInstance     string
	pathParams       url.Values
	method           string
}

// TODO: better test cases
var insertTests = []treeTestCase{
	{fakeHandler, "", "", url.Values{}, http.MethodGet},
	{fakeHandler, "/profile", "/profile", url.Values{}, http.MethodGet},
	{fakeHandler, "/profile/{id}", "/profile/2", url.Values{"id": []string{"2"}}, http.MethodGet},
	{fakeHandler, "/register", "/register", url.Values{}, http.MethodPost},
	{fakeHandler, "/user/profile/posts/{id}", "/user/profile/posts/3", url.Values{"id": []string{"3"}}, http.MethodPost},
	{fakeHandler, "/user/profile/posts/{id}", "/user/profile/posts/3", url.Values{"id": []string{"3"}}, http.MethodGet},
	{fakeHandler, "/user/profile/posts/{id}", "/user/profile/posts/3", url.Values{"id": []string{"3"}}, http.MethodDelete},
	{fakeHandler, "/user/profile/posts/{id}", "/user/profile/posts/3", url.Values{"id": []string{"3"}}, http.MethodPut},
	{fakeHandler, "/user/profile/posts/{id}", "/user/profile/posts/3", url.Values{"id": []string{"3"}}, http.MethodPatch},
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
	tree := node{}
	var ok bool
	for _, testCase := range insertTests {
		ok = tree.insert(
			parts(testCase.pathWithVariable),
			testCase.method,
			testCase.handler,
		)
		if !ok {
			t.Errorf("Failed to insert route %v", testCase)
		}
	}
	for _, testCase := range insertTests {
		handler, pathParams, statusCode := tree.findHandler(testCase.pathInstance, testCase.method)
		if fmt.Sprint(testCase.handler) != fmt.Sprint(handler) {
			t.Errorf("Tree find error: expected %v got %v", testCase.handler, handler)
		}

		if !reflect.DeepEqual(testCase.pathParams, pathParams) {
			t.Errorf("Tree find error: expected %v got %v", testCase.pathParams, pathParams)
		}

		if !reflect.DeepEqual(http.StatusOK, statusCode) {
			t.Errorf("Tree find error: expected %v got %v", http.StatusOK, statusCode)
		}
	}

	for _, testCase := range notFoundTests {
		handler, _, statusCode := tree.findHandler(testCase.path, testCase.method)
		if statusCode != http.StatusNotFound {
			t.Errorf("Tree find error: route %v expected %v got %v", testCase.path, http.StatusNotFound, statusCode)
		}

		if handler != nil {
			t.Errorf("Tree find error: expected %v got %v", nil, handler)
		}
	}

	for _, testCase := range methodNotAllowedTests {
		_, _, statusCode := tree.findHandler(testCase.path, testCase.method)
		if statusCode != http.StatusMethodNotAllowed {
			t.Errorf("Tree find error: expected %v got %v", http.StatusMethodNotAllowed, statusCode)
		}
	}
}
