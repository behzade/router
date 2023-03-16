package router

import "net/http"

type Middleware interface {
    Pipe(handler http.Handler) http.Handler
}
