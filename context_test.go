package router

import (
	"context"
	"net/url"
	"reflect"
	"testing"
)

func TestUrlParams(t *testing.T) {
	ctx := context.Background()
	pathParams := url.Values{
		"var1":       []string{"asd", "test"},
		"product-id": []string{"324"},
	}
	newCtx := setPathParams(ctx, pathParams)

	newParams := GetPathParams(newCtx)
	if !reflect.DeepEqual(pathParams, newParams) {
		t.Errorf("Url Params error: expected %q got %q", pathParams, newParams)
	}
}
