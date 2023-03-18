package router

import (
	"context"
	"net/url"
)

type urlParams struct {
    PathParams url.Values
    QueryParams url.Values
}

type urlParamsKey string

const defaultPathParamsKey urlParamsKey = "path:params"

func setUrlParams(ctx context.Context, pathParams url.Values, queryParams url.Values) context.Context {
	return context.WithValue(ctx, defaultPathParamsKey, &urlParams{pathParams, queryParams})
}

func GetPathParams(ctx context.Context) (urlParams, bool) {
	urlParams, ok := ctx.Value(defaultPathParamsKey).(*urlParams)
	return *urlParams, ok
}
