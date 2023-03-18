package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/behzade/router"
)

func main() {
	r := router.New()

	r.GET("/", &testHandler{})
	r.GET("/asd/{v1}/{v2}/{v1}/qqqq/asd/{v1}?asd=22", &testHandler{})

	for {
		http.ListenAndServe("127.0.0.1:8000", r)
	}

}

type testHandler struct{}

func (t *testHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte(req.URL.Path))
}

func testFunc(ctx context.Context, i input) output {
	params, ok := router.GetUrlParams(ctx)
	if !ok {
		panic("asd")
	}
	fmt.Printf("params: %v\n", params)

	return output{i.Name}
}

type input struct {
	Name string
}

type output struct {
	Name2 string
}
