package logger

import (
	"context"
	"errors"
	"io/fs"
	"syscall"
	"time"

	"github.com/Trendyol/chaki/util/appctx"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	loggerCtxValuer = appctx.NewValuer[*zap.Logger]("zap_logger_key").OnDefault(New)
	initialized     = false
)

const DefaultTimeKey = "timestamp"

func Init(timeEncoder zapcore.TimeEncoder, level string, timeKey string, location *time.Location) error {
	if initialized {
		return errors.New("logger already initliazed")
	}

	zc := zap.NewProductionConfig()

	if timeKey != "" {
		zc.EncoderConfig.TimeKey = timeKey
	} else {
		zc.EncoderConfig.TimeKey = DefaultTimeKey
	}

	if timeEncoder == nil {
		timeEncoder = zapcore.EpochTimeEncoder
	}

	if location != nil {
		timeEncoder = wrapTimeEncoderWithLocation(timeEncoder, location)
	}

	zc.EncoderConfig.EncodeTime = timeEncoder

	if level != "" {
		var zapLevel zapcore.Level
		if err := zapLevel.UnmarshalText([]byte(level)); err == nil {
			zc.Level = zap.NewAtomicLevelAt(zapLevel)
		}
	}

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
	err := New().Sync()

	var pathErr *fs.PathError
	if errors.Is(err, syscall.ENOTTY) || errors.As(err, &pathErr) {
		return nil
	}

	return err
}

func From(ctx context.Context) *zap.Logger {
	return loggerCtxValuer.Get(ctx)
}

func WithLogger(parent context.Context, logger *zap.Logger) context.Context {
	return loggerCtxValuer.Set(parent, logger)
}

func wrapTimeEncoderWithLocation(encoder zapcore.TimeEncoder, location *time.Location) zapcore.TimeEncoder {
	return func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		encoder(t.In(location), enc)
	}
}
