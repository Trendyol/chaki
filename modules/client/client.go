package client

import (
	"context"

	"github.com/Trendyol/chaki/config"
	"github.com/go-resty/resty/v2"
)

type Base struct {
	name    string
	driver  *resty.Client
	circuit *circuit
	rc      *retryConfig
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

	clientCfg := f.cfg.Of("client").Of(name)

	return &Base{
		name: name,
		driver: newDriverBuilder(clientCfg).
			AddErrDecoder(cOpts.errDecoder).
			AddUpdaters(f.baseWrappers...).
			AddUpdaters(cOpts.driverWrappers...).
			build(),
		circuit: newCircuit(clientCfg, name),
		rc:      getRetryConfigs(clientCfg),
	}
}

func (b *Base) Request(ctx context.Context) *Request {

	return &Request{
		circuit:     b.circuit,
		Request:     b.driver.R().SetContext(ctx),
		errF:        defaultCircuitErrorFunc,
		retryConfig: b.rc,
	}
}
