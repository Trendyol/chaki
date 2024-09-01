package otelkafka

import (
	"context"
	"errors"

	"github.com/Trendyol/chaki/module"
	"github.com/Trendyol/chaki/modules/kafka/consumer"
	"github.com/Trendyol/chaki/modules/kafka/producer"
	otelmodule "github.com/Trendyol/chaki/modules/otel"
	"go.opentelemetry.io/otel/trace"
)

func WithKafka() otelmodule.Option {
	m := module.NewSubModule().Provide(
		newSpanBuilder,
		newConsumerInterceptor,
		newProducerInterceptor,
		newBatchConsumerInterceptor,
	)

	return otelmodule.WithSubModule(m)
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

func newBatchConsumerInterceptor(sb *spanBuilder) consumer.BatchConsumerInterceptor {
	return consumer.BatchConsumerInterceptorFunc(func(msgs []*consumer.Message, next consumer.BatchConsumeFn) error {
		spans := make([]trace.Span, len(msgs))
		defer func() {
			for _, span := range spans {
				span.End()
			}
		}()

		for i, msg := range msgs {
			ctx, span := sb.BuildConsumerSpan(msg)
			msg.Context = ctx
			spans[i] = span
		}

		// Record single message errors
		for _, msg := range msgs {
			if ed := msg.ErrDescription; ed != "" {
				span := trace.SpanFromContext(msg.Context)
				span.RecordError(errors.New(ed))
			}
		}

		// Record all message errors
		if err := next(msgs); err != nil {
			for _, span := range spans {
				span.RecordError(err)
			}
			return err
		}

		return nil
	})
}

func newProducerInterceptor(sb *spanBuilder) producer.Interceptor {
	return producer.InterceptorFunc(func(ctx context.Context, msgs []producer.Message, next producer.InterceptNextFunc) error {
		ctx, span := sb.BuildProducerSpan(ctx, msgs)
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
