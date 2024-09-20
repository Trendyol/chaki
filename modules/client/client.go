package client

import (
	"context"

	"github.com/Trendyol/chaki/config"
	"github.com/go-resty/resty/v2"
)

type Base struct {
	name   string
	driver *resty.Client
}

type Factory struct {
	cfg          *config.Config
	baseWrappers []DriverWrapper
}

func NewFactory(cfg *config.Config, wrappers []DriverWrapper) *Factory {
	return &Factory{
		cfg:          cfg,
		baseWrappers: wrappers,
	}
}

func (f *Factory) Get(name string, opts ...Option) *Base {
	cOpts := &options{
		errDecoder:     DefaultErrDecoder(name),
		driverWrappers: []DriverWrapper{},
	}

	for _, opt := range opts {
		opt.Apply(cOpts)
	}

	return &Base{
		driver: newDriverBuilder(f.cfg.Of("client").Of(name)).
			AddErrDecoder(cOpts.errDecoder).
			AddUpdaters(f.baseWrappers...).
			AddUpdaters(cOpts.driverWrappers...).
			build(),
		name: name,
	}
}

func (r *Base) Request(ctx context.Context) *resty.Request {
	return r.driver.R().SetContext(ctx)
}
