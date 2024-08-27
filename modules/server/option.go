package server

import (
	"github.com/Trendyol/chaki/modules/server/middlewares"
	"github.com/gofiber/fiber/v2"
)

type options struct {
	errHandler fiber.ErrorHandler
}

type Option interface {
	Apply(*options)
}

type withFunc func(*options)

func (wf withFunc) Apply(opts *options) {
	wf(opts)
}

func WithErrorHandler(errHandler fiber.ErrorHandler) Option {
	return withFunc(func(o *options) {
		o.errHandler = errHandler
	})
}

func buildOptions(option ...Option) *options {
	opts := &options{
		errHandler: middlewares.ErrHandler(),
	}

	for _, o := range option {
		o.Apply(opts)
	}

	return opts
}
