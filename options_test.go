package router

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"
	"testing"
)

type testResponeWriter struct {
	writer bytes.Buffer
	header http.Header
}

func (w testResponeWriter) Header() http.Header {
	return w.header
}

func (w testResponeWriter) Write(data []byte) (int, error) {
	return w.writer.Write(data)
}

func (w testResponeWriter) WriteHeader(statusCode int) {
	w.header.Set("status", fmt.Sprint(statusCode))

}

func TestOptionsHandler(t *testing.T) {
	handler := OptionsHandler{[]string{http.MethodGet, http.MethodPost}, http.StatusOK}
	w := testResponeWriter{bytes.Buffer{}, http.Header{}}
	req := http.Request{Method: http.MethodOptions, URL: &url.URL{Path: "/"}}
	handler.ServeHTTP(w, &req)
	if w.header.Get("Allow") != "GET, POST" {
		t.Errorf("Options Handler error: expected %q got %q", "GET, POST", w.header.Get("Allow"))
	}
}
