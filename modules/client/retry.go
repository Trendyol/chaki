package client

import (
	"github.com/Trendyol/chaki/config"
	"time"
)

type retryCfg struct {
	Enabled  bool
	Attempts uint
	Delay    time.Duration
}

func retryFromCfg(config *config.Config) *retryCfg {
	if config == nil {
		return &retryCfg{
			Enabled: false,
		}
	}
	enabled := config.GetBool("enabled")
	attempts := config.GetUInt("attempts")
	delay := config.GetDuration("delay")

	return &retryCfg{
		Enabled:  enabled,
		Attempts: attempts,
		Delay:    delay,
	}
}
