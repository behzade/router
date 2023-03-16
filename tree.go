package router

import (
	"net/http"
)

type Tree struct {
	children map[string]*Tree  // NOTE: map key is path part
	handlers   map[string]http.Handler // NOTE: map key is http method
}

func (root *Tree) insert(splitPath []string, route http.Handler, method string) bool {
	if len(splitPath) == 0 {
		_, ok := root.handlers[method]
		if ok {
			return false
		}
		root.handlers[method] = route
		return true
	}

	child, ok := root.children[splitPath[0]]

	if ok {
		return child.insert(splitPath[1:], route, method)

	}

	child = &Tree{map[string]*Tree{}, map[string]http.Handler{}}
	child.insert(splitPath[1:], route, method)
	root.children[splitPath[0]] = child
	return true
}

func (root *Tree) find(splitPath []string, method string) (http.Handler, int) {
	if len(splitPath) == 0 {
		if len(root.handlers) == 0 {
			return nil, http.StatusNotFound
		}

		route, ok := root.handlers[method]
		if ok {
			return route, http.StatusOK
		}
		return nil, http.StatusMethodNotAllowed
	}

	child, ok := root.children[splitPath[0]]

	if ok {
		return child.find(splitPath[1:], method)
	}

	return nil, http.StatusNotFound
}
