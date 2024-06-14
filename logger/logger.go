package logger

import (
	"context"
	"errors"

	"github.com/Trendyol/chaki/util/appctx"
	"go.uber.org/zap"
)

var (
	loggerCtxValuer = appctx.NewValuer[*zap.Logger]("zap_logger_key").OnDefault(New)
	initialized     = false
)

func Init() error {
	if initialized {
		return errors.New("logger already initliazed")
	}

	zc := zap.NewProductionConfig()
	zc.EncoderConfig.TimeKey = "timestamp"

	logger, err := zc.Build()
	if err != nil {
		return err
	}

	zap.ReplaceGlobals(logger)

	initialized = true
	return nil
}

func Fatal(err error) {
	if err == nil {
		return
	}

	New().Fatal(err.Error())
}

func New() *zap.Logger {
	return zap.L()
}

func Sync() error {
	return New().Sync()
}

func From(ctx context.Context) *zap.Logger {
	return loggerCtxValuer.Get(ctx)
}

func WithLogger(parent context.Context, logger *zap.Logger) context.Context {
	return loggerCtxValuer.Set(parent, logger)
}
