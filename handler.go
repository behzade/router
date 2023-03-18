package router

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gorilla/schema"
)

type JsonHandler[T any, R any] func(context context.Context, requestData T) R

func ToHttpHandler[T any, R any](fn JsonHandler[T, R]) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		data := new(T)
		if req.Method == http.MethodGet {
			encoder := schema.NewDecoder()
			encoder.IgnoreUnknownKeys(true)
			err := encoder.Decode(data, req.URL.Query())
			if err != nil {
				panic(err)
			}
		} else {
			err := json.NewDecoder(req.Body).Decode(data)
			if err != nil {
				panic(err)
			}
		}

		result := fn(req.Context(), *data)
		json.NewEncoder(w).Encode(result)
	})
}
