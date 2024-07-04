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

func (f *Factory) Get(name string, errDecoder ErrDecoder, driverWrappers ...DriverWrapper) *Base {
	return &Base{
		driver: newDriverBuilder(f.cfg.Of("client").Of(name)).
			AddErrDecoder(errDecoder).
			AddUpdaters(f.baseWrappers...).
			AddUpdaters(driverWrappers...).
			build(),
		name: name,
	}
}

func (r *Base) Request(ctx context.Context) *resty.Request {
	return r.driver.R().SetContext(ctx)
}
