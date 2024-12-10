package client

import (
	"context"

	"github.com/Trendyol/chaki/config"
	"github.com/Trendyol/chaki/modules/client/common"
	"github.com/go-resty/resty/v2"
)

type Base struct {
	name   string
	driver *resty.Client
}

type Factory struct {
	cfg          *config.Config
	baseWrappers []DriverWrapper
	rtWrappers   []common.RoundTripperWrapper
}

func NewFactory(cfg *config.Config, wrappers []DriverWrapper, rtWrappers []common.RoundTripperWrapper) *Factory {
	initCircuitPresets(cfg)
	initRetryPresets(cfg)
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
			SetRetry(getRetryConfigs(clientCfg)).
			SetCircuit(getCircuitConfigs(clientCfg)).
			build(),
	}
}

func (b *Base) Request(ctx context.Context) *resty.Request {
	return b.driver.R().SetContext(ctx)
}

func (b *Base) RequestWithCommand(ctx context.Context, command string) *resty.Request {
	return b.driver.R().SetContext(context.WithValue(ctx, circuitCommandKey, command))
}
