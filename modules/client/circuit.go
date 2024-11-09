package client

import (
	"context"
	"github.com/Trendyol/chaki/config"
	"github.com/afex/hystrix-go/hystrix"
)

type (
	CircuitFunc        func(context.Context) error
	CircuitErrorFunc   func(context.Context, error) error
	CircuitErrorFilter func(error) (bool, error)

	circuitConfig struct {
		Name                   string
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
		Name: name,
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

func (c *circuit) do(ctx context.Context, fn CircuitFunc, fallback func(context.Context, error) error, fi ...CircuitErrorFilter) error {
	var e error
	var ok bool

	function := func(ctx context.Context) error {

		err := fn(ctx)

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
		return hystrixErr
	}

	if e != nil {
		return e
	}

	return nil
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
