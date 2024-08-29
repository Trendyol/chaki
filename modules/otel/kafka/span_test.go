package otelkafka

import (
	"context"
	"testing"

	"github.com/Trendyol/chaki/modules/common/ctxvaluer"
	"github.com/Trendyol/chaki/modules/kafka/consumer"
	"github.com/Trendyol/chaki/modules/kafka/producer"
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
	"go.opentelemetry.io/otel/trace/noop"
)

func buildTestSpanContext() context.Context {
	var (
		parentctx, _ = noop.NewTracerProvider().Tracer("test tracer").Start(context.Background(), "test span")
		traceID      = trace.TraceID([16]byte{0x01, 0x23, 0x45, 0x67, 0x89, 0xab, 0xcd, 0xef, 0xfe, 0xdc, 0xba, 0x98, 0x76, 0x54, 0x32, 0x10})
		spanID       = trace.SpanID([8]byte{0x00, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77})
	)

	parentctx = trace.ContextWithRemoteSpanContext(parentctx, trace.NewSpanContext(trace.SpanContextConfig{
		TraceID: traceID,
		SpanID:  spanID,
	}))

	return parentctx
}

func Test_spanBuilder_BuildConsumerSpan(t *testing.T) {
	// Given
	var (
		tp  = noop.NewTracerProvider()
		tmp = otel.GetTextMapPropagator()
		sb  = newSpanBuilder(tp, tmp)
		msg = &consumer.Message{
			Context: buildTestSpanContext(),
		}
	)

	// When
	ctx, span := sb.BuildConsumerSpan(msg)

	// Then
	assert.NotNil(t, span)
	assert.NotNil(t, ctxvaluer.SpanID.Get(ctx))
	assert.NotNil(t, ctxvaluer.TraceID.Get(ctx))
}

func Test_spanBuilder_BuildProducerSpan(t *testing.T) {
	// Given
	var (
		tp        = noop.NewTracerProvider()
		tmp       = propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{})
		sb        = newSpanBuilder(tp, tmp)
		msgs      = []producer.Message{{}, {}}
		parentctx = buildTestSpanContext()
	)

	// When
	ctx, span := sb.BuildProducerSpan(parentctx, msgs)

	// Then
	assert.NotNil(t, span)
	assert.NotNil(t, ctxvaluer.SpanID.Get(ctx))
	assert.NotNil(t, ctxvaluer.TraceID.Get(ctx))
	assert.Len(t, msgs, 2)
	assert.Len(t, msgs[0].Headers, 1)
	assert.Equal(t, msgs[0].Headers[0].Key, "traceparent")
}
