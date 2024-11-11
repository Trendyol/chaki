package client

import (
	"context"

	"github.com/Trendyol/chaki/config"
	"github.com/Trendyol/chaki/util/store"
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
		name   string
	}
)

var (
	defaultCircuitConfig = &circuitConfig{
		Name:                   "default",
		Timeout:                5000,
		MaxConcurrentRequests:  100,
		ErrorPercentThreshold:  50,
		RequestVolumeThreshold: 20,
		SleepWindow:            5000,
	}

	aggressiveCircuitConfig = &circuitConfig{
		Name:                   "aggressive",
		Timeout:                2000,
		MaxConcurrentRequests:  50,
		ErrorPercentThreshold:  25,
		RequestVolumeThreshold: 10,
		SleepWindow:            3000,
	}

	relaxedCircuitConfig = &circuitConfig{
		Name:                   "relaxed",
		Timeout:                10000,
		MaxConcurrentRequests:  200,
		ErrorPercentThreshold:  75,
		RequestVolumeThreshold: 40,
		SleepWindow:            7000,
	}

	circuitPresetMap = store.NewBucket[string, *circuitConfig](func(k string) *circuitConfig { return nil })
)

func newCircuit(cfg *config.Config, name string) *circuit {
	c := getCircuitConfigs(cfg)

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
		name:   name,
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

func initCircuitPresets(cfg *config.Config) {
	presets := []*circuitConfig{
		defaultCircuitConfig,
		aggressiveCircuitConfig,
		relaxedCircuitConfig,
	}

	for _, cc := range presets {
		circuitPresetMap.Set(cc.Name, cc)
	}

	userPresets, err := config.ToStruct[[]*circuitConfig](cfg, "client.circuitPresets")
	if err != nil {
		panic(err)
	}
	for _, cc := range userPresets {
		circuitPresetMap.Set(cc.Name, cc)
	}
}

func getCircuitConfigs(cfg *config.Config) *circuitConfig {
	if !cfg.GetBool("circuit.enabled") {
		return &circuitConfig{}
	}

	preset := cfg.GetString("circuit.preset")
	switch preset {
	case "custom":
		cc, err := config.ToStruct[*circuitConfig](cfg, "circuit")
		if err != nil {
			panic(err)
		}
		return cc
	default:
		if cc := circuitPresetMap.Get(preset); cc != nil {
			return cc
		}
		panic("unknown circuit breaker preset: " + preset)
	}
}
