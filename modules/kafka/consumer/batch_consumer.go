package consumer

import (
	kafkalib "github.com/Trendyol/kafka-konsumer/v2"
)

type (
	Messages       = []kafkalib.Message
	BatchConsumeFn = kafkalib.BatchConsumeFn

	BatchConsumer interface {
		Name() string
		Consume([]*Message) error
	}

	BatchConsumerInterceptor interface {
		Intercept(msg []*Message, next BatchConsumeFn) error
	}

	BatchConsumerInterceptorFunc func(msg []*Message, next BatchConsumeFn) error

	batchConsumerImpl struct {
		name string
		fn   BatchConsumeFn
	}
)

var (
	// type-check
	_ Interceptor = InterceptorFunc(nil)
)

func notImplementedBatchConsumeFn([]*kafkalib.Message) error {
	panic("consume method not implemented")
}

func NewBatchConsumer(name string) BatchConsumer {
	return &batchConsumerImpl{
		name: name,
		fn:   notImplementedBatchConsumeFn,
	}
}

func NewBatchConsumerFn(name string, fn BatchConsumeFn) BatchConsumer {
	return &batchConsumerImpl{name, fn}
}

func (c *batchConsumerImpl) Name() string {
	return c.name
}

func (c *batchConsumerImpl) Consume(msg []*Message) error {
	return c.fn(msg)
}

func (f BatchConsumerInterceptorFunc) Intercept(messages []*Message, next BatchConsumeFn) error {
	return f(messages, next)
}
