package router

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

type Tree struct {
	staticChildren  map[string]*Tree        // map key is path part
	dynamicChildren map[string]*Tree        // map key is path variable part
	handlers        map[string]http.Handler // map key is http method
}

// add a new path to the router, does nothing and returns false on duplicate path,method pair
func (t *Tree) insert(pathParts []PathPart, method string, handler http.Handler) bool {
	if len(pathParts) == 0 {
		_, ok := t.handlers[method]
		if ok {
			return false
		}
		t.handlers[method] = handler
		return true
	}

	var ok bool
	var child *Tree

	if pathParts[0].IsVariable {
		child, ok = t.dynamicChildren[pathParts[0].Value]
	} else {
		child, ok = t.staticChildren[pathParts[0].Value]
	}

	if ok {
		return child.insert(pathParts[1:], method, handler)
	}

	child = &Tree{map[string]*Tree{}, map[string]*Tree{}, map[string]http.Handler{}}
	child.insert(pathParts[1:], method, handler)
	if pathParts[0].IsVariable {
		t.dynamicChildren[pathParts[0].Value] = child
	} else {
		t.staticChildren[pathParts[0].Value] = child
	}

	return true
}

func (root *Tree) findNode(path string, offset int, params url.Values) (*Tree, bool) {
	part, offset := parse(path, offset)
	if part == "" {
		return root, true
	}

	var child *Tree
	var ok bool

	child, ok = root.staticChildren[part]

	if ok {
		return child.findNode(path, offset, params)
	}

	var key string
	for key, child = range root.dynamicChildren {
		child, ok = child.findNode(path, offset, params)

		if ok {
			params.Add(key, part)
			return child, true
		}
	}
	return nil, false
}

func (root *Tree) findHandler(path string, method string) (http.Handler, url.Values, int) {
	params := url.Values{}

	node, ok := root.findNode(path, 0, params)

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

func (t *Tree) String() string {
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
