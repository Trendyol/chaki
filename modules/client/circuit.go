package client

import (
	"context"

	"github.com/Trendyol/chaki/config"
	"github.com/afex/hystrix-go/hystrix"
	"github.com/go-resty/resty/v2"
)

type (
	CircuitFunc        func(context.Context) (*resty.Response, error)
	CircuitErrorFunc   func(context.Context, error) error
	CircuitErrorFilter func(error) (bool, error)

	circuitConfig struct {
		Name                   string
		enabled                bool
		Timeout                int
		MaxConcurrentRequests  int
		ErrorPercentThreshold  int
		RequestVolumeThreshold int
		SleepWindow            int
		Commands               []string
	}

	circuit struct {
		config *circuitConfig
	}
)

func newCircuit(cfg *config.Config, name string) *circuit {
	c := &circuitConfig{
		Name:    name,
		enabled: cfg.GetBool("circuit.enabled"),
	}

	if cfg.GetBool("circuit.enabled") {
		c, err := config.ToStruct[*circuitConfig](cfg, "circuit")
		c.Name = name
		if err != nil {
			panic("could not convert the circuit for client:" + name + ". check your configuration.")
		}
	}

	// TODO: Presets
	hystrixConfig := hystrix.CommandConfig{
		Timeout:                c.Timeout,
		MaxConcurrentRequests:  c.MaxConcurrentRequests,
		ErrorPercentThreshold:  c.ErrorPercentThreshold,
		RequestVolumeThreshold: c.RequestVolumeThreshold,
		SleepWindow:            c.SleepWindow,
	}

	// TODO: circuit per endpoint?
	hystrix.ConfigureCommand(name, hystrixConfig)

	return &circuit{
		config: c,
	}
}

func (c *circuit) do(ctx context.Context, fn CircuitFunc, fallback func(context.Context, error) error, fi ...CircuitErrorFilter) (*resty.Response, error) {
	if c.config == nil || !c.config.enabled {
		return fn(ctx)
	}

	var e error
	var ok bool
	var resp *resty.Response

	function := func(ctx context.Context) error {

		var err error
		resp, err = fn(ctx)

		for _, filter := range fi {
			if ok, e = filter(err); ok {
				return err
			}
		}

		if len(fi) > 0 {
			return nil
		}

		return err
	}

	hystrixErr := hystrix.DoC(ctx, c.config.Name, function, fallback)

	if hystrixErr != nil {
		return nil, hystrixErr
	}

	if e != nil {
		return nil, e
	}

	return resp, nil
}

func defaultCircuitErrorFunc(_ context.Context, err error) error {
	return err
}

func setDefaultCircuitConfigs(cfg *config.Config) {
	cfg.SetDefault("circuit.enabled", false)
	cfg.SetDefault("circuit.timeout", 5000)
	cfg.SetDefault("circuit.maxConcurrentRequests", 100)
	cfg.SetDefault("circuit.requestVolumeThreshold", 20)
	cfg.SetDefault("circuit.sleepWindow", 5000)
	cfg.SetDefault("circuit.errorPercentThreshold", 50)
}
