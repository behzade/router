package router

import (
	"net/http"
	"net/url"
	"sort"
)

type Middleware interface {
	Next(handler http.Handler) http.Handler
}

type Router struct {
	NotFoundHandler         http.Handler
	MethodNotAllowedHandler http.Handler

	tree             *Tree
	middlewares      []Middleware // global middlewares
	middlewareScores []int
}

func New() *Router {
	return &Router{
		tree: &Tree{make(map[string]*Tree), make(map[string]*Tree), make(map[string]http.Handler)},
	}
}

// resolve the handler for path. returns the handler, pathParams and status code
func (r *Router) resolve(path string, method string) (http.Handler, url.Values, int) {
	handler, pathParams, statusCode := r.tree.find(parse(path), method)

	switch statusCode {
	case http.StatusNotFound:
		return r.NotFoundHandler, nil, statusCode
	case http.StatusMethodNotAllowed:
		if r.MethodNotAllowedHandler != nil {
			handler = r.MethodNotAllowedHandler
		}
		return handler, pathParams, statusCode
	default:
		return handler, pathParams, statusCode
	}
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	handler, pathParams, statusCode := r.resolve(req.URL.Path, req.Method)

	if handler == nil {
		w.WriteHeader(statusCode)
		return
	}

	for _, middleware := range r.middlewares {
		handler = middleware.Next(handler)
	}

	handler.ServeHTTP(w, req.WithContext(setPathParams(req.Context(), pathParams)))
}

func (r *Router) addRoute(method string, path string, handler http.Handler) {
	r.tree.insert(parts(path), method, handler)
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

func (r *Router) AddMiddleware(score int, middleware Middleware) {
	index := sort.SearchInts(r.middlewareScores, score)
	insertToIndex(r.middlewareScores, index, score)
	insertToIndex(r.middlewares, index, middleware)
}
