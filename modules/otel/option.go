package otel

import "context"

type CloseFunc func(context.Context) error

type InitFunc func() CloseFunc

type options struct {
	initFunc InitFunc
}

type Option interface {
	Apply(*options) *options
}

type withOption func(*options) *options

func (wf withOption) Apply(opts *options) *options {
	return wf(opts)
}

var _ Option = withOption(nil)

func WithInitFunc(f InitFunc) Option {
	return withOption(func(o *options) *options {
		o.initFunc = f
		return o
	})
}

func applyOptions(opts ...Option) *options {
	o := &options{}

	for _, opt := range opts {
		o = opt.Apply(o)
	}

	return o
}
