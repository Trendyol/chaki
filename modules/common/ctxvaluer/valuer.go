package ctxvaluer

import (
	"context"
	"github.com/Trendyol/chaki/logger"
	"github.com/Trendyol/chaki/util/appctx"
	"go.uber.org/zap"
)

const (
	CorrelationIdKey = "x-correlationId"
	ExecutorUserKey  = "x-executor-user"
	TraceIdKey       = "trace-id"
	SpanIdKey        = "span-id"
	AgentNameKey     = "x-agent-name"
)

var (
	CorrelationId = appctx.NewValuer[string](CorrelationIdKey)
	ExecutorUser  = appctx.NewValuer[string](ExecutorUserKey)
	TraceId       = appctx.NewValuer[string](TraceIdKey)
	SpanId        = appctx.NewValuer[string](SpanIdKey)
	AgentName     = appctx.NewValuer[string](AgentNameKey)
)

type CreateParams struct {
	CorrelationId string
	ExecutorUser  string
	AgentName     string
}

func CreateBaseTaskContext(parent context.Context, params CreateParams) context.Context {
	ctx := parent
	ctx = CorrelationId.Set(ctx, params.CorrelationId)
	ctx = ExecutorUser.Set(ctx, params.ExecutorUser)
	ctx = AgentName.Set(ctx, params.AgentName)

	l := logger.New().With(
		zap.String(CorrelationIdKey, params.CorrelationId),
		zap.String(ExecutorUserKey, params.ExecutorUser),
		zap.String(AgentNameKey, params.AgentName),
	)

	return logger.WithLogger(ctx, l)
}
