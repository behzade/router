package router

import (
	"net/http"
	"strings"
)

type OptionsHandler struct {
	options []string
}

func (o *OptionsHandler) ServeHTTP(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Allow", strings.Join(o.options, ", "))
	w.WriteHeader(http.StatusNoContent)
}
