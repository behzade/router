package router

import (
	"net/http"
	"sort"
)

type Handler func(w http.ResponseWriter, req *http.Request, pathParams Params)

type Middleware interface {
	Next(handler Handler) Handler
}

type Router struct {
	NotFoundHandler         Handler
	MethodNotAllowedHandler Handler

	tree             *node
	middlewares      []Middleware // global middlewares
	middlewareScores []int
}

func New() *Router {
	return &Router{
		tree: &node{},
	}
}

// resolve the handler for path. returns the handler, pathParams and status code
func (r *Router) resolve(path string, method string) (Handler, Params, int) {
	handler, pathParams, statusCode := r.tree.findHandler(path, method)

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

	handler(w, req, pathParams)
}

func (r *Router) AddRoute(method string, path string, handler Handler) {
	r.tree.insert(parts(path), method, handler)
}

func (r *Router) GET(path string, handler Handler) {
	r.AddRoute(http.MethodGet, path, handler)
}

func (r *Router) POST(path string, handler Handler) {
	r.AddRoute(http.MethodPost, path, handler)
}

func (r *Router) PUT(path string, handler Handler) {
	r.AddRoute(http.MethodPut, path, handler)
}

func (r *Router) PATCH(path string, handler Handler) {
	r.AddRoute(http.MethodPatch, path, handler)
}

func (r *Router) DELETE(path string, handler Handler) {
	r.AddRoute(http.MethodDelete, path, handler)
}

func (r *Router) AddMiddleware(score int, middleware Middleware) {
	index := sort.SearchInts(r.middlewareScores, score)
	insertToIndex(r.middlewareScores, index, score)
	insertToIndex(r.middlewares, index, middleware)
}

func (r *Router) String() string {
    return r.tree.String()
}
