package router

import (
	"context"
	"net/url"
)

type urlParamsKey string

const defaultPathParamsKey urlParamsKey = "path:params"

func setPathParams(ctx context.Context, pathParams url.Values) context.Context {
	return context.WithValue(ctx, defaultPathParamsKey, pathParams)
}

func GetPathParams(ctx context.Context) url.Values {
	params, ok := ctx.Value(defaultPathParamsKey).(url.Values)
	if !ok {
		panic("invalid context value")
	}
	return params
}
