package router

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gorilla/schema"
)

type Handler[T any, R any] func(requestData T, context context.Context) R

func ToHttpHandler[T any, R any](fn Handler[T, R]) http.HandlerFunc{
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
        data := new(T)
		if req.Method == http.MethodGet {
			err := schema.NewDecoder().Decode(data, req.URL.Query())
			if err != nil {
				panic(err)
			}
		} else {
			err := json.NewDecoder(req.Body).Decode(data)
			if err != nil {
				panic(err)
			}
		}

		result := fn(*data, req.Context())
		json.NewEncoder(w).Encode(result)
	})
}
