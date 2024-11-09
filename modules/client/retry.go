package client

import (
	"github.com/Trendyol/chaki/config"
	"time"
)

const (
	ConstantDelay    = "constant"
	IncrementalDelay = "incremental"
)

// TODO: Presets
type retryConfig struct {
	Count      int           `json:"count"`
	Interval   time.Duration `json:"interval"`
	MaxDelay   time.Duration `json:"maxDelay"`
	Multiplier float64       `json:"multiplier"`
	DelayType  string        `json:"delayType"`
}

func getRetryConfigs(cfg *config.Config) *retryConfig {
	if !cfg.GetBool("retry.enabled") {
		return &retryConfig{}
	}

	rc, err := config.ToStruct[*retryConfig](cfg, "retry")
	if err != nil {
		panic(err)
	}

	return rc
}

func setDefaultRetryConfigs(cfg *config.Config) {
	cfg.SetDefault("retry.enabled", false)
	cfg.SetDefault("retry.count", 3)
	cfg.SetDefault("retry.interval", "100ms")
	cfg.SetDefault("retry.maxDelay", "5s")
	cfg.SetDefault("retry.multiplier", 1.0)
	cfg.SetDefault("retry.delayType", "constant")
}
