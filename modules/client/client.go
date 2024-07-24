package client

import (
	"context"
	"fmt"
	"github.com/Trendyol/chaki/config"
	"github.com/go-resty/resty/v2"
)

type Base struct {
	name    string
	cfg     *config.Config
	driver  *resty.Client
	circuit *circuit
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
	cfg := f.cfg.Of("client").Of(name)

	driver := newDriverBuilder(cfg).
		AddErrDecoder(errDecoder).
		AddUpdaters(f.baseWrappers...).
		AddUpdaters(driverWrappers...).
		build()

	var circuit *circuit
	if cfg.Exists("circuit") {
		circuitCfg, err := config.ToStruct[circuitConfig](cfg, "circuit")
		if err != nil {
			panic("could not convert the circuit for client:" + name + ". check your configuration.")
		}

		circuit = newCircuit(&circuitCfg)
	}

	b := &Base{
		driver:  driver,
		name:    name,
		circuit: circuit,
		cfg:     cfg,
	}

	return b
}

func (b *Base) Request(ctx context.Context, name ...string) *Request {
	n := ""
	var rc *retryConfig
	if len(name) == 1 {
		n = name[0]
		rc = b.getRetryConfigs(n)
	}

	return &Request{
		name:        n,
		circuit:     b.circuit,
		Request:     b.driver.R().SetContext(ctx),
		errF:        defaultCircuitErrorFunc,
		retryConfig: rc,
	}
}

func (b *Base) getRetryConfigs(name string) *retryConfig {
	if b.cfg.Exists("requests") {
		cfg := b.cfg.Of("requests")
		setDefaultRetryConfigs(cfg, name)

		retryConfig, err := config.ToStruct[retryConfig](b.cfg.Of("requests"), name)
		if err != nil {
			panic(err)
		}

		return &retryConfig
	}

	panic(fmt.Sprintf("trying to initialize retry configs for request: %s but cannot find 'requests' tag on clients config", name))
}

func setDefaultRetryConfigs(c *config.Config, name string) {
	cfg := c.Of(name)
	cfg.SetDefault("retryCount", 3)
	cfg.SetDefault("initialDelay", "100ms")
	cfg.SetDefault("maxDelay", "100ms")
	cfg.SetDefault("multiplier", 1)
	cfg.SetDefault("delayType", "constant")
}
