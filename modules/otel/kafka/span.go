package otelkafka

import (
	"context"
	"fmt"
	"strconv"

	"github.com/Trendyol/chaki/modules/kafka/consumer"
	"github.com/Trendyol/chaki/modules/kafka/producer"
	"github.com/Trendyol/chaki/modules/otel/common"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"go.opentelemetry.io/otel/trace"
)

const tracerName = "github.com/Trendyol/chaki/otel/kafka"

type spanBuilder struct {
	tmp propagation.TextMapPropagator
	tp  trace.TracerProvider
}

func newSpanBuilder(tp trace.TracerProvider, tmp propagation.TextMapPropagator) *spanBuilder {
	return &spanBuilder{tmp, tp}
}

func (sb *spanBuilder) BuildConsumerSpan(msg *consumer.Message) (context.Context, trace.Span) {
	return sb.buildMessageSpan(msg, trace.SpanKindConsumer)
}

func (sb *spanBuilder) BuildProducerSpan(parent context.Context, msgs []producer.Message) (context.Context, trace.Span) {
	tracer := sb.tp.Tracer(tracerName)

	ctx, span := tracer.Start(parent, "producing message", trace.WithSpanKind(trace.SpanKindProducer))

	for i := range msgs {
		carrier := NewMessageCarrier(&msgs[i])
		sb.tmp.Inject(ctx, carrier)
	}

	return ctx, span
}

func (sb *spanBuilder) buildMessageSpan(msg *consumer.Message, kind trace.SpanKind) (context.Context, trace.Span) {
	var (
		tracer        = sb.tp.Tracer(tracerName)
		carrier       = NewMessageCarrier(msg)
		parentSpanCtx = sb.tmp.Extract(msg.Context, carrier)
	)

	attrs := []attribute.KeyValue{
		semconv.MessagingSystem("kafka"),
		semconv.MessagingDestinationKindTopic,
		semconv.MessagingDestinationName(msg.Topic),
		semconv.MessagingOperationReceive,
		semconv.MessagingMessageID(strconv.FormatInt(msg.Offset, 10)),
		semconv.MessagingKafkaSourcePartition(int(msg.Partition)),
	}

	opts := []trace.SpanStartOption{
		trace.WithAttributes(attrs...),
		trace.WithSpanKind(kind),
	}

	ctx, span := tracer.Start(parentSpanCtx, fmt.Sprintf("%s receive", msg.Topic), opts...)

	// Inject current span context, so consumers can use it to propagate span.
	sb.tmp.Inject(ctx, carrier)

	return common.WithContext(ctx, span), span
}
