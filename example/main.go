package main

import (
	"fmt"
	"net/http"

	"github.com/behzade/router"
)

func main() {
	r := router.New()

	r.GET("/", &indexHandler{})
	r.GET("/v1/product/{product-id}/comments",
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			productId := router.GetPathParams(r.Context()).Get("product-id")

			fmt.Fprintf(w, "Product %v comments", productId)
		}))

	r.POST("/v1/product/{product-id}/comments",
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Do stuff
		}))

	for {
		http.ListenAndServe("127.0.0.1:8000", r)
	}

}

type indexHandler struct {
}

func (t *indexHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("hello world!"))
}
