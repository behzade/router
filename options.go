package router

import (
	"net/http"
	"strings"
)

type OptionsHandler struct {
	options    []string
	statusCode int
}

func (o OptionsHandler) ServeHTTP(w http.ResponseWriter, _ *http.Request, _ Params) {
	w.Header().Set("Allow", strings.Join(o.options, ", "))
	w.WriteHeader(o.statusCode)
}
