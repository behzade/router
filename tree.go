package router

import (
	"net/http"
)

type Tree struct {
	staticChildren  map[string]*Tree        // NOTE: map key is path part
	dynamicChildren map[string]*Tree        // NOTE: map key is path variable part
	handlers        map[string]http.Handler // NOTE: map key is http method
}

func (t *Tree) insert(pathParts []PathPart, route http.Handler, method string) bool {
	if len(pathParts) == 0 {
		_, ok := t.handlers[method]
		if ok {
			return false
		}
		t.handlers[method] = route
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
		return child.insert(pathParts[1:], route, method)
	}

	child = &Tree{map[string]*Tree{}, map[string]*Tree{}, map[string]http.Handler{}}
	child.insert(pathParts[1:], route, method)
	if pathParts[0].IsVariable {
		t.dynamicChildren[pathParts[0].Value] = child
	} else {
		t.staticChildren[pathParts[0].Value] = child
	}

	return true
}

func (t *Tree) find(pathParts []string, method string) (http.Handler, map[string]string, int) {
    defaultStatus := http.StatusNotFound
	if len(pathParts) == 0 {
		if len(t.handlers) == 0 {
			return nil, nil, http.StatusNotFound
		}

		route, ok := t.handlers[method]
		if ok {
			return route, map[string]string{}, http.StatusOK
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
			pathParams[key] = pathParts[0]
			return handler, pathParams, status
		}

        if status == http.StatusMethodNotAllowed {
            defaultStatus = http.StatusMethodNotAllowed
        }
	}

	return nil, nil, defaultStatus
}
