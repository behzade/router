package router

import (
	"context"
	"net/url"
	"reflect"
	"testing"
)

func TestUrlParams(t *testing.T) {
	ctx := context.Background()
	params := urlParams{
		url.Values{
			"var1":       []string{"asd", "test"},
			"product-id": []string{"324"},
		},
		url.Values{
			"id":   []string{"1"},
			"sort": []string{"asc"},
		},
	}
	newCtx := setUrlParams(ctx, params.PathParams, params.QueryParams)

	newParams, ok := GetUrlParams(newCtx)
	if !ok {
		t.Error("Url Params error: failed to get")
	}
	if !reflect.DeepEqual(params, newParams) {
		t.Errorf("Url Params error: expected %q got %q", params, newParams)
	}
}
