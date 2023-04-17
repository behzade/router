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
	handlers        map[string]Handler
	pathParts       []pathPart
}

// add a new path to the router, does nothing and returns false on duplicate path,method pair
func (t *node) insert(pathParts []pathPart, method string, handler Handler) bool {
	if len(pathParts) == 0 {
		if t.handlers == nil {
			t.handlers = map[string]Handler{}
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

func (root *node) findNode(path string) (*node, url.Values) {
	if path == "" {
		return root, nil
	}

	part, rest := parse(path)
	if len(part) == 0 {
		return root, nil
	}

	if child, ok := root.staticChildren[string(part)]; ok {
		return child.findNode(rest)
	}

    var child *node
	var params url.Values
	var key string

	for key, child = range root.dynamicChildren {
		child, params = child.findNode(rest)

		if child != nil {
			if params == nil {
				params = url.Values{}
			}
			params.Add(key, string(part))
			return child, params
		}
	}
	return nil, nil
}

func (root *node) findHandler(path string, method string) (Handler, url.Values, int) {
	node, params := root.findNode(path)

	if node == nil {
		return nil, nil, http.StatusNotFound
	}

	if method == http.MethodOptions {
		handler := OptionsHandler{keys(root.handlers), http.StatusOK}
		return handler.ServeHTTP, nil, http.StatusOK
	}

	if len(node.handlers) == 0 {
		return nil, nil, http.StatusNotFound
	}

	handler, ok := node.handlers[method]

	if !ok {
		handler := &OptionsHandler{keys(root.handlers), http.StatusMethodNotAllowed}
		return handler.ServeHTTP, params, http.StatusMethodNotAllowed
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
