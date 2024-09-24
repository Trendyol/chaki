package ctxvaluer

import (
	"context"

	"github.com/Trendyol/chaki/logger"
	"github.com/Trendyol/chaki/util/appctx"
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

type CreateParams struct {
	CorrelationID string
	ExecutorUser  string
	AgentName     string
	Owner         string
}

func CreateBaseTaskContext(parent context.Context, params CreateParams) context.Context {
	ctx := parent
	ctx = CorrelationID.Set(ctx, params.CorrelationID)
	ctx = ExecutorUser.Set(ctx, params.ExecutorUser)
	ctx = AgentName.Set(ctx, params.AgentName)
	ctx = Owner.Set(ctx, params.Owner)

	l := logger.New().With(
		zap.String(CorrelationIDKey, params.CorrelationID),
		zap.String(ExecutorUserKey, params.ExecutorUser),
		zap.String(AgentNameKey, params.AgentName),
		zap.String(OwnerKey, params.Owner),
	)

	return logger.WithLogger(ctx, l)
}
