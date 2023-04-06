package router

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

type node struct {
	staticChildren  map[string]*node
	dynamicChildren map[string]*node
	handlers        map[string]http.Handler
	pathParts       []pathPart
}

// add a new path to the router, does nothing and returns false on duplicate path,method pair
func (t *node) insert(pathParts []pathPart, method string, handler http.Handler) bool {
	if len(pathParts) == 0 {
		if t.handlers == nil {
			t.handlers = map[string]http.Handler{}
		}
		_, ok := t.handlers[method]
		if ok {
			return false
		}
		t.handlers[method] = handler
		return true
	}

	var ok bool
	var child *node

	if pathParts[0].IsVariable && t.dynamicChildren != nil {
		child, ok = t.dynamicChildren[pathParts[0].Value]
	} else if t.staticChildren != nil {
		child, ok = t.staticChildren[pathParts[0].Value]
	}

	if ok {
		return child.insert(pathParts[1:], method, handler)
	}

	child = &node{}
	child.insert(pathParts[1:], method, handler)
	if pathParts[0].IsVariable {
		if t.dynamicChildren == nil {
			t.dynamicChildren = map[string]*node{pathParts[0].Value: child}
		} else {
			t.dynamicChildren[pathParts[0].Value] = child
		}
	} else {
		if t.staticChildren == nil {
			t.staticChildren = map[string]*node{pathParts[0].Value: child}
		} else {
			t.staticChildren[pathParts[0].Value] = child
		}
	}

	return true
}

func (root *node) findNode(path string, params url.Values) (*node, bool) {
	if path == "" {
		return root, true
	}

	part, rest := parse(path)
	if len(part) == 0 {
		return root, true
	}

	var child *node
	var ok bool

	child, ok = root.staticChildren[string(part)]

	if ok {
		return child.findNode(rest, params)
	}

	var key string
	for key, child = range root.dynamicChildren {
		child, ok = child.findNode(rest, params)

		if ok {
			params.Add(key, string(part))
			return child, true
		}
	}
	return nil, false
}

func (root *node) findHandler(path string, method string) (http.Handler, url.Values, int) {
	params := url.Values{}

	node, ok := root.findNode(path, params)

	if !ok {
		return nil, nil, http.StatusNotFound
	}
	if method == http.MethodOptions {
		return &OptionsHandler{keys(root.handlers), http.StatusOK}, nil, http.StatusOK
	}

	if len(node.handlers) == 0 {
		return nil, params, http.StatusNotFound
	}

	var handler http.Handler

	handler, ok = node.handlers[method]

	if !ok {
		return &OptionsHandler{keys(root.handlers), http.StatusMethodNotAllowed}, params, http.StatusMethodNotAllowed
	}
	return handler, params, http.StatusOK
}

func (t *node) String() string {
	var builder strings.Builder
	if len(t.handlers) != 0 {
		builder.WriteString(fmt.Sprintf(": %v\n", keys(t.handlers)))
	}

	for route, child := range t.staticChildren {
		builder.WriteString(fmt.Sprintf("%v/%v", route, child.String()))
	}

	for route, child := range t.dynamicChildren {
		builder.WriteString(fmt.Sprintf("{%v}/%v", route, child.String()))
	}
	return builder.String()
}
