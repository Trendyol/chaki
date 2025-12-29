package logger

import (
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

	return Init(timeEncoder, level, timeKey)
}
