package router

import (
	"net/http"
	"strings"
)

type OptionsHandler struct {
	options []string
}

func (o *OptionsHandler) ServeHTTP(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Allow", o.allowedString())
	w.WriteHeader(http.StatusNoContent)
}

func (o *OptionsHandler) allowedString() string {
    return strings.Join(o.options, ", ")
}
