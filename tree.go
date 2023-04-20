package router

import (
	"fmt"
	"net/http"
	"strings"
)

type namedNode struct {
	*node
	name string
}

var maxDynamicCount uint8 = 0

type node struct {
	staticChildren map[string]*node
	dynamicChild   *namedNode
	handlers       map[string]Handler
	pathParts      []pathPart
}

// add a new path to the router, does nothing and returns false on duplicate path,method pair
func (t *node) insert(pathParts pathParts, method string, handler Handler) bool {
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

	if pathParts[0].IsVariable && t.dynamicChild != nil {
		if pathParts[0].Value != t.dynamicChild.name {
			return false
		}
		child = t.dynamicChild.node
        ok = true;
	} else if t.staticChildren != nil {
		child, ok = t.staticChildren[pathParts[0].Value]
	}

	if ok {
		return child.insert(pathParts[1:], method, handler)
	}

	child = &node{}
	child.insert(pathParts[1:], method, handler)
	if pathParts[0].IsVariable {
		t.dynamicChild = &namedNode{child, pathParts[0].Value}
	} else {
		if t.staticChildren == nil {
			t.staticChildren = map[string]*node{pathParts[0].Value: child}
		} else {
			t.staticChildren[pathParts[0].Value] = child
		}
	}

    maxDynamicCount = max(uint8(pathParts.dynamicCount()), maxDynamicCount)
	return true
}

var buf [256]byte

func (root *node) findNode(path string) (*node, Params) {
	if path == "" {
		return root, nil
	}

	var n int
	var i int

	for ; i < len(path); i++ {
		c := path[i]
		if c == '/' && n > 0 {
			i++
			break
		}

		if c >= 'a' && c <= 'z' || c >= '0' && c <= '9' || c == '-' {
			buf[n] = c
			n++
		} else if c >= 'A' && c <= 'Z' {
			buf[n] = c + 32 // to lower
			n++
		}

	}

	if n == 0 {
		return root, nil
	}

	if child, ok := root.staticChildren[string(buf[:n])]; ok {
		return child.findNode(path[i:])
	}

	var child *node
	var params Params

	if root.dynamicChild != nil {
		child, params = root.dynamicChild.findNode(path[i:])

		if child != nil {
			if params == nil {
				params = make(Params, 0, maxDynamicCount)
			}
            length := len(params)
            params = params[:length + 1]
            params[length].key = root.dynamicChild.name
            params[length].value = buf[:n]
			return child, params
		}
	}

	return nil, nil
}

func (root *node) findHandler(path string, method string) (Handler, Params, int) {
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

	if t.dynamicChild != nil {
		builder.WriteString(fmt.Sprintf("{%v}/%v", t.dynamicChild.name, t.dynamicChild.String()))
	}
	return builder.String()
}
