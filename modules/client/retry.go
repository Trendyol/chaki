package client

import (
	"time"

	"github.com/Trendyol/chaki/config"
	"github.com/Trendyol/chaki/util/store"
)

type DelayType string

const (
	ConstantDelay    DelayType = "constant"
	ExponentialDelay DelayType = "exponential"
)

type retryConfig struct {
	Name      string        `json:"name"`
	Count     int           `json:"count"`
	Interval  time.Duration `json:"interval"`
	MaxDelay  time.Duration `json:"maxDelay"`
	DelayType DelayType     `json:"delayType"`
}

var retryPresetMap = store.NewBucket[string, *retryConfig](func(k string) *retryConfig { return nil })

var (
	defaultRetryConfig = &retryConfig{
		Name:      "default",
		Count:     3,
		Interval:  100 * time.Millisecond,
		MaxDelay:  5 * time.Second,
		DelayType: ConstantDelay,
	}
	exponentialRetryConfig = &retryConfig{
		Name:      "exponential",
		Count:     3,
		Interval:  100 * time.Millisecond,
		MaxDelay:  5 * time.Second,
		DelayType: ExponentialDelay,
	}
	aggresiveRetryConfig = &retryConfig{
		Name:      "aggresive",
		Count:     7,
		Interval:  50 * time.Millisecond,
		MaxDelay:  2 * time.Second,
		DelayType: ConstantDelay,
	}
	aggresiveExponentialRetryConfig = &retryConfig{
		Name:      "aggresiveExponential",
		Count:     7,
		Interval:  50 * time.Millisecond,
		MaxDelay:  2 * time.Second,
		DelayType: ExponentialDelay,
	}
	relaxedRetryConfig = &retryConfig{
		Name:      "relaxed",
		Count:     2,
		Interval:  500 * time.Millisecond,
		MaxDelay:  2 * time.Second,
		DelayType: ConstantDelay,
	}
	relaxedExponentialConfig = &retryConfig{
		Name:      "relaxedExponential",
		Count:     2,
		Interval:  500 * time.Millisecond,
		MaxDelay:  2 * time.Second,
		DelayType: ExponentialDelay,
	}

	predefinedRetryPresets = []*retryConfig{
		defaultRetryConfig,
		exponentialRetryConfig,
		aggresiveRetryConfig,
		aggresiveExponentialRetryConfig,
		relaxedRetryConfig,
		relaxedExponentialConfig,
	}
)

func getRetryConfigs(cfg *config.Config) *retryConfig {
	if !cfg.GetBool("retry.enabled") {
		return &retryConfig{}
	}

	preset := cfg.GetString("retry.preset")
	switch preset {
	case "custom":
		rc, err := config.ToStruct[*retryConfig](cfg, "retry")
		if err != nil {
			panic(err)
		}
		return rc
	default:
		if rc := retryPresetMap.Get(preset); rc != nil {
			return rc
		}

		panic("unknown retry preset: " + preset)
	}
}

func setDefaultRetryConfigs(cfg *config.Config) {
	cfg.SetDefault("retry.enabled", false)
	cfg.SetDefault("retry.preset", "default")
	cfg.SetDefault("retry.count", 3)
	cfg.SetDefault("retry.interval", "100ms")
	cfg.SetDefault("retry.maxDelay", "5s")
	cfg.SetDefault("retry.delayType", "constant")
}

func initRetryPresets(cfg *config.Config) {
	for _, rc := range predefinedRetryPresets {
		retryPresetMap.Set(rc.Name, rc)
	}

	userPresets, err := config.ToStruct[[]*retryConfig](cfg, "client.retryPresets")
	if err != nil {
		panic(err)
	}
	for _, rc := range userPresets {
		retryPresetMap.Set(rc.Name, rc)
	}
}
