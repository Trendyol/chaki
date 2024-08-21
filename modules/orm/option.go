package orm

import "github.com/Trendyol/chaki/modules/orm/driver"

type options struct {
	driver driver.Driver
}

type Option interface {
	Apply(*options) *options
}

func buildOptions(opts ...Option) *options {
	o := &options{}

	for _, opt := range opts {
		o = opt.Apply(o)
	}

	return o
}
