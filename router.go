package router

import (
	"net/http"
	"sort"
)

type Router struct {
	NotFoundHandler         http.Handler
	MethodNotAllowedHandler http.Handler

	handlers         *Tree
	middlewares      []Middleware // global middlewares
	middlewareScores []int
}

func New() *Router {
	return &Router{
		handlers: &Tree{make(map[string]*Tree), make(map[string]*Tree), make(map[string]http.Handler)},
	}
}

func (r *Router) resolve(path string, method string) (http.Handler, map[string]string, int) {
	route, pathParams, statusCode := r.handlers.find(split(path), method)

	switch statusCode {
	case http.StatusOK:
		return route, pathParams, statusCode
	case http.StatusNotFound:
		return r.NotFoundHandler, nil, statusCode
	case http.StatusMethodNotAllowed:
		return r.MethodNotAllowedHandler, nil, statusCode
	}

	return r.NotFoundHandler, nil, http.StatusNotFound
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	handler, _, statusCode := r.resolve(req.URL.Path, req.Method)

	if statusCode != http.StatusOK {
		if handler != nil {
			handler.ServeHTTP(w, req)
            return
		}
		w.WriteHeader(statusCode)
		return
	}

	for _, middleware := range r.middlewares {
		handler = middleware.Pipe(handler)
	}

	handler.ServeHTTP(w, req)
}

func (r *Router) addRoute(method string, path string, handler http.Handler) {
	r.handlers.insert(parts(path), handler, method)
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

func (r *Router) AddMiddleware(score int, middleware Middleware) {
	index := sort.SearchInts(r.middlewareScores, score)
	insertToIndex(r.middlewareScores, index, score)
	insertToIndex(r.middlewares, index, middleware)
}
