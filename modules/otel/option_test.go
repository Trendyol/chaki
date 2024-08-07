package otel

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWithInitFunc(t *testing.T) {
	// Given
	initf := func() CloseFunc {
		return func(ctx context.Context) error { return nil }
	}
	opt := WithInitFunc(initf)

	// When
	opts := applyOptions(opt)

	// Then
	assert.NotNil(t, opts.initFunc)
}
