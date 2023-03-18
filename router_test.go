package router

import (
	"bytes"
	"net/http"
	"testing"
)

// taken from httprouter
func TestRouter(t *testing.T) {
	var get, post, put, patch, delete bool

	router := New()
	router.GET("/GET", http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		get = true
	}))
	router.POST("/POST", http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		post = true
	}))
	router.PUT("/PUT", http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		put = true
	}))
	router.PATCH("/PATCH", http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		patch = true
	}))
	router.DELETE("/DELETE", http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		delete = true
	}))

	w := testResponeWriter{bytes.Buffer{}, http.Header{}}

	r, _ := http.NewRequest(http.MethodGet, "/GET", nil)
	router.ServeHTTP(w, r)
	if !get {
		t.Error("routing GET failed")
	}

	r, _ = http.NewRequest(http.MethodPost, "/POST", nil)
	router.ServeHTTP(w, r)
	if !post {
		t.Error("routing POST failed")
	}

	r, _ = http.NewRequest(http.MethodPut, "/PUT", nil)
	router.ServeHTTP(w, r)
	if !put {
		t.Error("routing PUT failed")
	}

	r, _ = http.NewRequest(http.MethodPatch, "/PATCH", nil)
	router.ServeHTTP(w, r)
	if !patch {
		t.Error("routing PATCH failed")
	}

	r, _ = http.NewRequest(http.MethodDelete, "/DELETE", nil)
	router.ServeHTTP(w, r)
	if !delete {
		t.Error("routing DELETE failed")
	}
}
