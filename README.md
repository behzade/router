# Simple Router
This is a simple router that is compatible with the standard library http.Handler interface.

## Features
- Remove extra slashes: multiple slashes or trailing slashes have no effect on the routing
- Handle route parameters: add named paramters to the path with {var} syntax. Path params are captured and saved in the context as a value.
- Builtin Options handler: get the Allow header with an OPTIONS request or on MethodNotAllowed responses.
- Can use custom not found and method not allowed handlers.
- Add global middlewares to the router with a score to be sorted by.

## Limitations
- Only supports lowercase a-z and 0-9 and "-" as allowed characters in path.

## Example Usage
```go
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
```

## Issues
Feel free to open an issue with a bug/feature request.
