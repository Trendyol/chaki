package logger

import (
	"github.com/Trendyol/chaki/config"
	"go.uber.org/zap/zapcore"
)

const (
	TimeEncoderEpoch       = "epoch"
	TimeEncoderISO8601     = "iso8601"
	TimeEncoderRFC3339     = "rfc3339"
	TimeEncoderRFC3339Nano = "rfc3339nano"

	DefaultTimezone = "Local"
)

var defaultConfigMap = map[string]any{
	"timeEncoder": TimeEncoderEpoch,
	"level":       "info",
	"timeKey":     DefaultTimeKey,
	"timezone":    DefaultTimezone,
}

var timeEncoderMap = map[string]zapcore.TimeEncoder{
	TimeEncoderEpoch:       zapcore.EpochTimeEncoder,
	TimeEncoderISO8601:     zapcore.ISO8601TimeEncoder,
	TimeEncoderRFC3339:     zapcore.RFC3339TimeEncoder,
	TimeEncoderRFC3339Nano: zapcore.RFC3339NanoTimeEncoder,
}

func setDefaultConfigs(cfg *config.Config) {
	for k, v := range defaultConfigMap {
		cfg.SetDefault(k, v)
	}
}

func getTimeEncoder(name string) zapcore.TimeEncoder {
	if encoder, ok := timeEncoderMap[name]; ok {
		return encoder
	}
	return zapcore.EpochTimeEncoder
}
