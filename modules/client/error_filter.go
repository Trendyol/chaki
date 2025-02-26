package client

import (
	"context"
)

// WIP
type errorFilterFunc func(error) (bool, error)

func SetErrorFilter(ctx context.Context, filter errorFilterFunc) context.Context {
	return context.WithValue(ctx, circuitErrFilterKey, filter)
}

func getErrorFilterFunc(ctx context.Context) errorFilterFunc {
	if fb, ok := ctx.Value(circuitErrFilterKey).(errorFilterFunc); ok {
		return fb
	}
	return defaultErrorFilter
}

func defaultErrorFilter(err error) (bool, error) {
	return true, err
}
