package otelkafka

import (
	"context"
	"testing"

	"github.com/Trendyol/chaki/modules/common/ctxvaluer"
	"github.com/Trendyol/chaki/modules/kafka/consumer"
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace/noop"
)

func Test_spanBuilder_BuildConsumerSpan(t *testing.T) {
	// Given
	var (
		tp  = noop.NewTracerProvider()
		tmp = otel.GetTextMapPropagator()
		sb  = newSpanBuilder(tp, tmp)
		msg = &consumer.Message{
			Context: context.Background(),
		}
	)

	// When
	ctx, span := sb.BuildConsumerSpan(msg)

	// Then
	assert.NotNil(t, span)
	assert.NotNil(t, ctxvaluer.SpanId.Get(ctx))
	assert.NotNil(t, ctxvaluer.TraceId.Get(ctx))

}

func Test_spanBuilder_BuildProducerSpan(t *testing.T) {
	// Given
	var (
		tp  = noop.NewTracerProvider()
		tmp = otel.GetTextMapPropagator()
		sb  = newSpanBuilder(tp, tmp)
		msg = &consumer.Message{
			Context: context.Background(),
		}
	)

	// When
	ctx, span := sb.BuildProducerSpan(msg)

	// Then
	assert.NotNil(t, span)
	assert.NotNil(t, ctxvaluer.SpanId.Get(ctx))
	assert.NotNil(t, ctxvaluer.TraceId.Get(ctx))

}
