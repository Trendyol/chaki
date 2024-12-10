package client

import (
	"context"
	"fmt"

	"github.com/Trendyol/chaki/config"
	"github.com/Trendyol/chaki/util/store"
	"github.com/afex/hystrix-go/hystrix"
)

type (
	circuitConfig struct {
		Name                   string
		Enabled                bool
		Timeout                int
		MaxConcurrentRequests  int
		ErrorPercentThreshold  int
		RequestVolumeThreshold int
		SleepWindow            int
	}

	contextKey string
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

	circuitPresetMap = store.NewBucket(func(k string) *circuitConfig { return nil })
)

const (
	circuitCommandKey   contextKey = "command"
	circuitFallbackKey  contextKey = "fallback"
	circuitErrFilterKey contextKey = "errorFilter"
)

func SetFallbackFunc(ctx context.Context, fb func(context.Context, error) (interface{}, error)) context.Context {
	return context.WithValue(ctx, circuitFallbackKey, fb)
}

func SetErrorFilter(ctx context.Context, filter func(error) (bool, error)) context.Context {
	return context.WithValue(ctx, circuitErrFilterKey, filter)
}

func defaultCircuitErrorFunc(commandName string) func(_ context.Context, err error) error {
	return func(_ context.Context, err error) error {
		return fmt.Errorf("command %s, error: %w", commandName, err)
	}
}

func setDefaultCircuitConfigs(cfg *config.Config) {
	cfg.SetDefault("circuit.enabled", false)
	cfg.SetDefault("circuit.preset", "default")
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
		return nil
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

func (c *circuitConfig) toHystrixConfig() hystrix.CommandConfig {
	return hystrix.CommandConfig{
		Timeout:                c.Timeout,
		MaxConcurrentRequests:  c.MaxConcurrentRequests,
		ErrorPercentThreshold:  c.ErrorPercentThreshold,
		RequestVolumeThreshold: c.RequestVolumeThreshold,
		SleepWindow:            c.SleepWindow,
	}
}
