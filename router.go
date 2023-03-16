package router

import "net/http"

type Router struct {
	NotFoundHandler         http.Handler
	MethodNotAllowedHandler http.Handler

	routeTree               *Tree
	middleware              []Middleware
}

func New() *Router {
	return &Router{}
}

func (r *Router) addRoute(path string, method string, handler http.Handler) {
	splitPath, _ := ParsePath(path)

	r.routeTree.insert(splitPath, &Route{handler, method})
}

func (r *Router) resolve(path string, method string) (http.Handler, bool) {
	splitPath, _ := ParsePath(path)
	route, ok := r.routeTree.find(splitPath)

	if !ok {
		return r.NotFoundHandler, false
	}

	if route.method != method {
		return r.MethodNotAllowedHandler, false
	}

	return route.handler, true
}
