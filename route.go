package router

import "net/http"

type Route struct {
	handler http.Handler
	name    *string
}

type Tree struct {
	children map[string]*Tree
	routes   map[string]*Route
}

func (root *Tree) insert(splitPath []string, route *Route, method string) {
	if len(splitPath) == 0 {
		root.routes[method] = route
		return
	}

	child, ok := root.children[splitPath[0]]

	if ok {
		child.insert(splitPath[1:], route, method)
		return
	}

	child = &Tree{map[string]*Tree{}, nil}
	child.insert(splitPath[1:], route, method)
	root.children[splitPath[0]] = child
}

func (root *Tree) find(splitPath []string, method string) (*Route, bool) {
	if len(splitPath) == 0 {
		route, ok := root.routes[method]
		return route, ok
	}

	child, ok := root.children[splitPath[0]]

	if ok {
		return child.find(splitPath[1:], method)
	}

	return nil, false
}
