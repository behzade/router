package router

import (
	"net/http"
)

type Router struct {
	NotFoundHandler         http.Handler
	MethodNotAllowedHandler http.Handler

	routes      *Tree
	middlewares []Middleware
}

func New() *Router {
	return &Router{
		NotFoundHandler:         http.HandlerFunc(NotFoundHandlerFunc),
		MethodNotAllowedHandler: http.HandlerFunc(MethodNotAllowedHandlerFunc),
		routes:                  &Tree{make(map[string]*Tree), make(map[string]*Route)},
	}
}

func (r *Router) addRoute(method string, path string, handler http.Handler) {
	splitPath := SplitPath(path)

	r.routes.insert(splitPath, &Route{handler}, method)
}

func (r *Router) GET(path string, handler http.Handler) {
	r.addRoute(http.MethodGet, path, handler)
}

func (r *Router) POST(path string, handler http.Handler) {
	r.addRoute(http.MethodPost, path, handler)
}

func (r *Router) PUT(path string, handler http.Handler) {
	r.addRoute(http.MethodPut, path, handler)
}

func (r *Router) PATCH(path string, handler http.Handler) {
	r.addRoute(http.MethodPatch, path, handler)
}

func (r *Router) DELETE(path string, handler http.Handler) {
	r.addRoute(http.MethodDelete, path, handler)
}

func (r *Router) HEAD(path string, handler http.Handler) {
	r.addRoute(http.MethodHead, path, handler)
}

func (r *Router) resolve(path string, method string) (http.Handler, int) {
	splitPath := SplitPath(path)
	route, statusCode := r.routes.find(splitPath, method)

	switch statusCode {
	case http.StatusOK:
		return route.handler, statusCode
	case http.StatusNotFound:
		return r.NotFoundHandler, statusCode
	case http.StatusMethodNotAllowed:
		return r.MethodNotAllowedHandler, statusCode
	}

	return r.NotFoundHandler, http.StatusNotFound
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	handler, statusCode := r.resolve(req.URL.Path, req.Method)
	if statusCode != http.StatusOK {
		handler.ServeHTTP(w, req)
		return
	}

	for _, middleware := range r.middlewares {
		handler = middleware(handler)
	}
	handler.ServeHTTP(w, req)
}
