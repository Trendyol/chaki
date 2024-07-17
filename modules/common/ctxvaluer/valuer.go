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
	OwnerKey         = "x-owner"
)

var (
	CorrelationId = appctx.NewValuer[string](CorrelationIdKey)
	ExecutorUser  = appctx.NewValuer[string](ExecutorUserKey)
	TraceId       = appctx.NewValuer[string](TraceIdKey)
	SpanId        = appctx.NewValuer[string](SpanIdKey)
	AgentName     = appctx.NewValuer[string](AgentNameKey)
	Owner         = appctx.NewValuer[string](OwnerKey)
)

type CreateParams struct {
	CorrelationId string
	ExecutorUser  string
	AgentName     string
	Owner         string
}

func CreateBaseTaskContext(parent context.Context, params CreateParams) context.Context {
	ctx := parent
	ctx = CorrelationId.Set(ctx, params.CorrelationId)
	ctx = ExecutorUser.Set(ctx, params.ExecutorUser)
	ctx = AgentName.Set(ctx, params.AgentName)
	ctx = Owner.Set(ctx, params.Owner)

	l := logger.New().With(
		zap.String(CorrelationIdKey, params.CorrelationId),
		zap.String(ExecutorUserKey, params.ExecutorUser),
		zap.String(AgentNameKey, params.AgentName),
		zap.String(OwnerKey, params.Owner),
	)

	return logger.WithLogger(ctx, l)
}
