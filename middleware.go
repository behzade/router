package router

import "net/http"

type Middleware interface {
    Next(handler http.Handler) http.Handler
}
