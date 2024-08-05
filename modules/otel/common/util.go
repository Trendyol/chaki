package common

import (
	"context"

	"github.com/Trendyol/chaki/logger"
	"github.com/Trendyol/chaki/modules/common/ctxvaluer"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

func WithContext(ctx context.Context, span trace.Span) context.Context {
	tid := span.SpanContext().TraceID().String()
	sid := span.SpanContext().SpanID().String()

	ctx = ctxvaluer.TraceId.Set(ctx, tid)
	ctx = ctxvaluer.SpanId.Set(ctx, sid)
	ctx = logger.WithLogger(ctx, logger.From(ctx).With(
		zap.String(ctxvaluer.TraceIdKey, tid),
		zap.String(ctxvaluer.SpanIdKey, sid),
	))

	return ctx
}
