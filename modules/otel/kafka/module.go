package otelkafka

import (
	"context"

	"github.com/Trendyol/chaki/module"
	"github.com/Trendyol/chaki/modules/kafka/consumer"
	"github.com/Trendyol/chaki/modules/kafka/producer"
	"go.opentelemetry.io/otel/trace"
)

const moduleName = "chaki-otel-kafka"

func Module() *module.Module {
	m := module.New(moduleName)

	m.Provide(
		newSpanBuilder,
		newConsumerInterceptor,
		newProducerInterceptor,
	)

	return m
}

func newConsumerInterceptor(sb *spanBuilder) consumer.Interceptor {
	return consumer.InterceptorFunc(func(msg *consumer.Message, next consumer.ConsumeFn) error {
		ctx, span := sb.BuildConsumerSpan(msg)
		defer span.End()

		msg.Context = ctx
		err := next(msg)

		if err != nil {
			span.RecordError(err)
			return err
		}

		return nil
	})
}

func newProducerInterceptor(sb *spanBuilder) producer.Interceptor {
	return producer.InterceptorFunc(func(ctx context.Context, msgs []producer.Message, next producer.InterceptNextFunc) error {
		var span trace.Span
		// TODO: update for bulk message producing
		if len(msgs) == 1 {
			msg := &msgs[0]
			msg.Context = ctx
			ctx, span = sb.BuildProducerSpan(msg)
		}

		err := next(ctx, msgs...)
		if err != nil && span != nil {
			span.RecordError(err)
		}

		if span != nil {
			span.End()
		}

		return err
	})
}
