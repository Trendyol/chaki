package ctxvaluer

import (
	"context"

	"github.com/Trendyol/chaki/logger"
	"github.com/Trendyol/chaki/util/appctx"
	"github.com/spf13/cast"
	"go.uber.org/zap"
)

const (
	CorrelationIDKey = "x-correlationId"
	ExecutorUserKey  = "x-executor-user"
	AgentNameKey     = "x-agentname"
	OwnerKey         = "x-owner"

	TraceIDKey = "trace-id"
	SpanIDKey  = "span-id"
)

var (
	CorrelationID = appctx.NewValuer[string](CorrelationIDKey)
	ExecutorUser  = appctx.NewValuer[string](ExecutorUserKey)
	TraceID       = appctx.NewValuer[string](TraceIDKey)
	SpanID        = appctx.NewValuer[string](SpanIDKey)
	AgentName     = appctx.NewValuer[string](AgentNameKey)
	Owner         = appctx.NewValuer[string](OwnerKey)
)

type CreateParams map[string]string

type stringValuer interface {
	Set(ctx context.Context, v string) context.Context
}

var valuers = map[string]stringValuer{
	CorrelationIDKey: CorrelationID,
	ExecutorUserKey:  ExecutorUser,
	AgentNameKey:     AgentName,
	OwnerKey:         Owner,
	TraceIDKey:       TraceID,
	SpanIDKey:        SpanID,
}

func CreateBaseTaskContext(parent context.Context, params CreateParams) context.Context {
	ctx := parent
	fields := make([]zap.Field, 0, len(params))

	for key, val := range params {
		if v, ok := valuers[key]; ok {
			ctx = v.Set(ctx, val)
		}
		fields = append(fields, zap.String(key, val))
	}

	l := logger.New().With(fields...)

	return logger.WithLogger(ctx, l)
}

func GetHeaderMapping(customHeaders map[string]any) map[string]string {
	mapping := map[string]string{
		CorrelationIDKey: CorrelationIDKey,
		ExecutorUserKey:  ExecutorUserKey,
		AgentNameKey:     AgentNameKey,
		OwnerKey:         OwnerKey,
	}

	for k, v := range customHeaders {
		mapping[k] = cast.ToString(v)
	}

	return mapping
}
