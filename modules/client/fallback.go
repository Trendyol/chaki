package client

import (
	"context"
	"net/http"
)

type fallbackFunc func(context.Context, error) (any, error)

func SetFallbackFunc(ctx context.Context, fb fallbackFunc) context.Context {
	return context.WithValue(ctx, circuitFallbackKey, fb)
}

type fallbackHandler struct {
	ctx      context.Context
	fn       func(context.Context, error) (interface{}, error)
	resp     *http.Response
	executed bool
}

func newFallbackHandler(ctx context.Context) *fallbackHandler {
	h := &fallbackHandler{
		ctx: ctx,
	}

	if fn, ok := ctx.Value(circuitFallbackKey).(fallbackFunc); ok {
		h.fn = fn
	} else {
		h.fn = defaultCircuitFallbackFunc
	}

	return h
}

func (f *fallbackHandler) handle(ctx context.Context, err error) error {
	resp, err := f.fn(ctx, err)
	if err != nil {
		return err
	}

	body, contentLength, contentType, err := interfaceToReadCloserWithLength(resp)
	if err != nil {
		return err
	}

	f.resp = &http.Response{
		StatusCode:    http.StatusOK,
		Status:        "200 OK",
		Body:          body,
		Header:        make(http.Header),
		ContentLength: contentLength,
	}

	if contentType != "" {
		f.resp.Header.Set("Content-Type", contentType)
	} else {
		f.resp.Header.Set("Content-Type", "application/json")
	}

	f.executed = true
	return nil
}

func defaultCircuitFallbackFunc(_ context.Context, err error) (any, error) {
	return nil, err
}
