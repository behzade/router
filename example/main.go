package main

import (
	"context"
	"net/http"

	"github.com/behzade/router"
)

func main() {
	r := router.New()

	r.GET("/", &testHandler{})
	r.GET("/asd/{var1}/{var2}/{var3}/qqqq/asd/{v1}?asd=22", &testHandler{})
	r.POST("/new", router.ToHttpHandler(testFunc))

	for {
		http.ListenAndServe("127.0.0.1:8000", r)
	}

}

type testHandler struct{}

func (t *testHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte(req.URL.Path))
}

func testFunc(_ context.Context, i input) output {
	return output{i.Name}
}

type input struct {
	Name string
}

type output struct {
	Name2 string
}
