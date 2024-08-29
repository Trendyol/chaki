package otelkafka

import (
	"context"
	"errors"
	"testing"

	"github.com/Trendyol/chaki/modules/common/ctxvaluer"
	"github.com/Trendyol/chaki/modules/kafka/consumer"
	"github.com/Trendyol/chaki/modules/kafka/producer"
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace/noop"
)

func Test_newConsumerInterceptor(t *testing.T) {
	t.Run("it should intercept and set context with otel information", func(t *testing.T) {
		// Given
		var (
			tp          = noop.NewTracerProvider()
			tmp         = otel.GetTextMapPropagator()
			sb          = newSpanBuilder(tp, tmp)
			interceptor = newConsumerInterceptor(sb)
			msg         = &consumer.Message{
				Context: context.Background(),
			}
		)

		nextf := func(msg *consumer.Message) error {
			assert.NotNil(t, ctxvaluer.TraceID.Get(msg.Context))
			assert.NotNil(t, ctxvaluer.SpanID.Get(msg.Context))
			return nil
		}

		// When
		err := interceptor.Intercept(msg, nextf)

		// Then
		assert.Nil(t, err)
	})

	t.Run("it should intercept and record err", func(t *testing.T) {
		// Given
		var (
			tp          = noop.NewTracerProvider()
			tmp         = otel.GetTextMapPropagator()
			sb          = newSpanBuilder(tp, tmp)
			interceptor = newConsumerInterceptor(sb)
			errTest     = errors.New("test-error")
			msg         = &consumer.Message{
				Context: context.Background(),
			}
		)

		nextf := func(msg *consumer.Message) error {
			return errTest
		}

		// When
		err := interceptor.Intercept(msg, nextf)

		// Then
		assert.Equal(t, errTest, err)
	})
}

func Test_newProducerInterceptor(t *testing.T) {
	t.Run("it should intercept and set context with otel information", func(t *testing.T) {
		// Given
		var (
			tp          = noop.NewTracerProvider()
			tmp         = otel.GetTextMapPropagator()
			sb          = newSpanBuilder(tp, tmp)
			interceptor = newProducerInterceptor(sb)
			msg         = producer.Message{
				Context: context.Background(),
			}
		)

		var nextf producer.InterceptNextFunc = func(ctx context.Context, msg ...producer.Message) error {
			assert.NotNil(t, ctxvaluer.TraceID.Get(msg[0].Context))
			assert.NotNil(t, ctxvaluer.SpanID.Get(msg[0].Context))
			return nil
		}

		// When
		err := interceptor.Intercept(msg.Context, []producer.Message{msg}, nextf)

		// Then
		assert.Nil(t, err)
	})

	t.Run("it should intercept and record err", func(t *testing.T) {
		// Given
		var (
			tp          = noop.NewTracerProvider()
			tmp         = otel.GetTextMapPropagator()
			sb          = newSpanBuilder(tp, tmp)
			interceptor = newProducerInterceptor(sb)
			errTest     = errors.New("test-error")
			msg         = producer.Message{
				Context: context.Background(),
			}
		)

		var nextf producer.InterceptNextFunc = func(ctx context.Context, msg ...producer.Message) error {
			assert.NotNil(t, ctxvaluer.TraceID.Get(msg[0].Context))
			assert.NotNil(t, ctxvaluer.SpanID.Get(msg[0].Context))
			return errTest
		}

		// When
		err := interceptor.Intercept(msg.Context, []producer.Message{msg}, nextf)

		// Then
		assert.Equal(t, errTest, err)
	})
}

func Test_newBatchConsumerInterceptor(t *testing.T) {
	t.Run("it should intercept and set context with otel information", func(t *testing.T) {
		// Given
		var (
			tp          = noop.NewTracerProvider()
			tmp         = otel.GetTextMapPropagator()
			sb          = newSpanBuilder(tp, tmp)
			interceptor = newBatchConsumerInterceptor(sb)
			msgs        = []*consumer.Message{
				{Context: context.Background()},
				{Context: context.Background()},
				{Context: context.Background()},
			}
		)

		var nextf consumer.BatchConsumeFn = func(msgs []*consumer.Message) error {
			for _, msg := range msgs {
				assert.NotNil(t, ctxvaluer.TraceID.Get(msg.Context))
				assert.NotNil(t, ctxvaluer.SpanID.Get(msg.Context))
			}
			return nil
		}

		// When
		err := interceptor.Intercept(msgs, nextf)

		// Then
		assert.Nil(t, err)
	})

	t.Run("it should intercept and record err", func(t *testing.T) {
		// Given
		var (
			tp          = noop.NewTracerProvider()
			tmp         = otel.GetTextMapPropagator()
			sb          = newSpanBuilder(tp, tmp)
			interceptor = newProducerInterceptor(sb)
			errTest     = errors.New("test-error")
			msg         = producer.Message{
				Context: context.Background(),
			}
		)

		var nextf producer.InterceptNextFunc = func(ctx context.Context, msg ...producer.Message) error {
			assert.NotNil(t, ctxvaluer.TraceID.Get(msg[0].Context))
			assert.NotNil(t, ctxvaluer.SpanID.Get(msg[0].Context))
			return errTest
		}

		// When
		err := interceptor.Intercept(msg.Context, []producer.Message{msg}, nextf)

		// Then
		assert.Equal(t, errTest, err)
	})
}
