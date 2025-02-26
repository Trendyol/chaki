package client

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
)

type fallbackFunc func(context.Context, error) (any, error)

func SetFallbackFunc(ctx context.Context, fb fallbackFunc) context.Context {
	return context.WithValue(ctx, circuitFallbackKey, fb)
}

type fallbackHandler struct {
	fn       func(context.Context, error) (interface{}, error)
	resp     *http.Response
	executed bool
}

func newOrDefaultFallbackHandler(ctx context.Context) *fallbackHandler {
	h := &fallbackHandler{}

	if fn, ok := ctx.Value(circuitFallbackKey).(fallbackFunc); ok {
		h.fn = fn
	} else {
		h.fn = defaultFallbackFn
	}

	return h
}

func defaultFallbackFn(_ context.Context, e error) (interface{}, error) {
	return nil, e
}

func (f *fallbackHandler) handle(ctx context.Context, e error) error {
	resp, err := f.fn(ctx, e)
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

func interfaceToReadCloserWithLength(data interface{}) (io.ReadCloser, int64, string, error) {
	b, err := json.Marshal(data)
	if err != nil {
		return nil, 0, "", err
	}
	return io.NopCloser(bytes.NewReader(b)), int64(len(b)), "application/json", nil
}
