package router

import (
	"net/http"
)

func NotFoundHandlerFunc(w http.ResponseWriter, _ *http.Request) {
    w.WriteHeader(http.StatusNotFound)
    w.Write([]byte("not found"))
}

func MethodNotAllowedHandlerFunc(w http.ResponseWriter, _ *http.Request) {
    w.WriteHeader(http.StatusMethodNotAllowed)
    w.Write([]byte("method not allowed"))
}
