package router

import (
	"bytes"
	"net/http"
	"net/url"
	"testing"
)

// taken from httprouter
func TestRouter(t *testing.T) {
	var get, post, put, patch, delete bool

	router := New()
	router.GET("/GET", func(http.ResponseWriter, *http.Request, url.Values) {
		get = true
	})
	router.POST("/POST", func(http.ResponseWriter, *http.Request, url.Values) {
		post = true
	})
	router.PUT("/PUT", func(http.ResponseWriter, *http.Request, url.Values) {
		put = true
	})
	router.PATCH("/PATCH", func(http.ResponseWriter, *http.Request, url.Values) {
		patch = true
	})
	router.DELETE("/DELETE", func(http.ResponseWriter, *http.Request, url.Values) {
		delete = true
	})

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
