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

func (t *Tree) findNode(pathParts []string, params url.Values) (*Tree, bool) {
	if len(pathParts) == 0 {
		return t, true
	}

	var child *Tree
	var ok bool

	child, ok = t.staticChildren[pathParts[0]]

	if ok {
		return child.findNode(pathParts[1:], params)
	}

	var key string
	for key, child = range t.dynamicChildren {
		child, ok = child.findNode(pathParts[1:], params)

		if ok {
			params.Add(key, pathParts[0])
			return child, true
		}
	}
	return nil, false
}

func (t *Tree) findHandler(pathParts []string, method string) (http.Handler, url.Values, int) {
	var node *Tree
	var ok bool

	params := url.Values{}

	node, ok = t.findNode(pathParts, params)

	if !ok {
		return nil, nil, http.StatusNotFound
	}
	if method == http.MethodOptions {
		return &OptionsHandler{keys(t.handlers), http.StatusOK}, nil, http.StatusOK
	}

	if len(node.handlers) == 0 {
		return nil, params, http.StatusNotFound
	}

	var handler http.Handler

	handler, ok = node.handlers[method]

	if !ok {
		return &OptionsHandler{keys(t.handlers), http.StatusMethodNotAllowed}, params, http.StatusMethodNotAllowed
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
