package otel

import (
	"context"

	"github.com/Trendyol/chaki/module"
)

type CloseFunc func(context.Context) error

type InitFunc func() CloseFunc

type options struct {
	initFunc   InitFunc
	subModules []*module.SubModule
}

type Option interface {
	Apply(*options)
}

type withOption func(*options)

func (wf withOption) Apply(opts *options) {
	wf(opts)
}

var _ Option = withOption(nil)

func WithInitFunc(f InitFunc) Option {
	return withOption(func(o *options) {
		o.initFunc = f
	})
}

func WithSubModule(sub ...*module.SubModule) Option {
	return withOption(func(o *options) {
		o.subModules = append(o.subModules, sub...)
	})
}

func applyOptions(opts ...Option) *options {
	o := &options{}

	for _, opt := range opts {
		opt.Apply(o)
	}

	return o
}
