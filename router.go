package router

import (
	"net/http"
	"net/url"
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

func (r *Router) resolve(path string, method string) (http.Handler, url.Values, url.Values, int) {
	splitPath, queryParams := parse(path)
	route, pathParams, statusCode := r.handlers.find(splitPath, method)

	switch statusCode {
	case http.StatusNotFound:
		return r.NotFoundHandler, nil, nil, statusCode
	case http.StatusMethodNotAllowed:
		return r.MethodNotAllowedHandler, nil, nil, statusCode
	default:
		return route, pathParams, queryParams, statusCode
	}
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	handler, pathParams, queryParams, statusCode := r.resolve(req.URL.Path, req.Method)

	if statusCode != http.StatusOK && handler == nil {
		w.WriteHeader(statusCode)
		return
	}

	for _, middleware := range r.middlewares {
		handler = middleware.Next(handler)
	}

	handler.ServeHTTP(w, req.WithContext(setUrlParams(req.Context(), pathParams, queryParams)))
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
