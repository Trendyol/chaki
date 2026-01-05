package logger

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Trendyol/chaki/config"
	"github.com/Trendyol/chaki/module"
	"go.uber.org/fx"
	"go.uber.org/zap/zapcore"
)

func Module() *module.Module {
	m := module.New("logger")

	m.Invoke(initLogger)

	return m
}

func initLogger(lc fx.Lifecycle, cfg *config.Config) {
	lc.Append(fx.StartStopHook(
		func() error {
			return initLoggerFromConfig(cfg)
		},
		Sync,
	))
}

func initLoggerFromConfig(cfg *config.Config) error {
	loggerCfg := cfg.Of("logger")
	setDefaultConfigs(loggerCfg)

	var timeEncoder zapcore.TimeEncoder
	timeEncoderName := loggerCfg.GetString("timeEncoder")
	timeEncoder = getTimeEncoder(timeEncoderName)

	level := loggerCfg.GetString("level")
	timeKey := loggerCfg.GetString("timeKey")
	timezone := loggerCfg.GetString("timezone")

	location, err := parseLocation(timezone)
	if err != nil {
		return err
	}

	return Init(timeEncoder, level, timeKey, location)
}

func parseLocation(timezone string) (*time.Location, error) {
	if timezone == "" || strings.EqualFold(timezone, "Local") {
		return time.Local, nil
	}

	if strings.EqualFold(timezone, "UTC") {
		return time.UTC, nil
	}

	if strings.HasPrefix(timezone, "UTC") {
		offsetStr := timezone[3:]
		if offsetStr == "" {
			return time.UTC, nil
		}

		offset, err := strconv.Atoi(offsetStr)
		if err != nil {
			return nil, fmt.Errorf("invalid timezone offset: %s", timezone)
		}

		return time.FixedZone(timezone, offset*60*60), nil
	}

	return time.LoadLocation(timezone)
}
