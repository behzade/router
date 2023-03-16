package router

import (
	"net/http"
)

type Tree struct {
	children map[string]*Tree  // NOTE: map key is path part
	handlers   map[string]http.Handler // NOTE: map key is http method
}

func (t *Tree) insert(pathParts []string, route http.Handler, method string) bool {
	if len(pathParts) == 0 {
		_, ok := t.handlers[method]
		if ok {
			return false
		}
		t.handlers[method] = route
		return true
	}

	child, ok := t.children[pathParts[0]]

	if ok {
		return child.insert(pathParts[1:], route, method)
	}

	child = &Tree{map[string]*Tree{}, map[string]http.Handler{}}
	child.insert(pathParts[1:], route, method)
	t.children[pathParts[0]] = child
	return true
}

func (t *Tree) find(pathParts []string, method string) (http.Handler, int) {
	if len(pathParts) == 0 {
		if len(t.handlers) == 0 {
			return nil, http.StatusNotFound
		}

		route, ok := t.handlers[method]
		if ok {
			return route, http.StatusOK
		}
		return nil, http.StatusMethodNotAllowed
	}

	child, ok := t.children[pathParts[0]]

	if ok {
		return child.find(pathParts[1:], method)
	}

	return nil, http.StatusNotFound
}
