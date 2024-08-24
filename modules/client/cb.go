package client

import (
	"github.com/Trendyol/chaki/config"
	"github.com/sony/gobreaker"
)

type cbCfg struct {
	Enabled bool
	Cb      *gobreaker.CircuitBreaker
}

func cbFromCfg(config *config.Config) *cbCfg {
	if config == nil {
		return &cbCfg{
			Enabled: false,
		}
	}

	cbSettings := gobreaker.Settings{
		Name:        config.GetString("name"),
		MaxRequests: config.GetUInt32("maxRequests"),
		Interval:    config.GetDuration("interval"),
		Timeout:     config.GetDuration("timeout"),
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			return counts.ConsecutiveFailures > config.GetUInt32("consecutiveFailures")
		},
	}
	return &cbCfg{
		Enabled: config.GetBool("enabled"),
		Cb:      gobreaker.NewCircuitBreaker(cbSettings),
	}

}
