package router

import (
	"net/http"
	"net/url"
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

func (t *Tree) find(pathParts []string, method string) (http.Handler, url.Values, int) {
	defaultStatus := http.StatusNotFound
	if len(pathParts) == 0 {
		if len(t.handlers) == 0 {
			return nil, nil, http.StatusNotFound
		}

		if method == http.MethodOptions {
			return &OptionsHandler{keys(t.handlers)}, url.Values{}, http.StatusNoContent
		}

		route, ok := t.handlers[method]
		if ok {
			return route, url.Values{}, http.StatusOK
		}
		return nil, nil, http.StatusMethodNotAllowed
	}

	child, ok := t.staticChildren[pathParts[0]]

	if ok {
		return child.find(pathParts[1:], method)
	}

	for key, child := range t.dynamicChildren {
		handler, pathParams, status := child.find(pathParts[1:], method)

		if status == http.StatusOK {
			pathParams.Add(key, pathParts[0])
			return handler, pathParams, status
		}

		if status == http.StatusMethodNotAllowed {
			defaultStatus = http.StatusMethodNotAllowed
		}
	}

	return nil, nil, defaultStatus
}
