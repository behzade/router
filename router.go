package router

import "net/http"

type Router struct {
	NotFoundHandler         http.Handler
	MethodNotAllowedHandler http.Handler

	routeTree   *Tree
	middlewares []Middleware
}

func New() *Router {
	return &Router{}
}

func (r *Router) addRoute(path string, name *string, method string, handler http.Handler) {
	splitPath := SplitPath(path)

	r.routeTree.insert(splitPath, &Route{handler}, method)
}

func (r *Router) resolve(path string, method string) (http.Handler, int) {
	splitPath := SplitPath(path)
	route, statusCode := r.routeTree.find(splitPath, method)

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
	}

	for _, middleware := range r.middlewares {
		handler = middleware(handler)
	}
	handler.ServeHTTP(w, req)
}
