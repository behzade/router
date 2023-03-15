package router

import "net/http"

type Route struct {
    handler http.Handler
}
