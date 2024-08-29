package consumer

import (
	kafkalib "github.com/Trendyol/kafka-konsumer/v2"
)

type (
	Message   = kafkalib.Message
	ConsumeFn = kafkalib.ConsumeFn

	// Consumer represent consumer func with name, so consumer configuration can be matched by using name
	Consumer interface {
		Name() string
		Consume(*Message) error
	}

	// Interceptor represent a middleware for consume functions
	Interceptor interface {
		Intercept(msg *Message, next ConsumeFn) error
	}

	// InterceptorFunc is a functional implementation of Interceptor
	InterceptorFunc func(msg *Message, next ConsumeFn) error

	consumerImpl struct {
		name string
		fn   ConsumeFn
	}
)

// type-check
var _ Interceptor = InterceptorFunc(nil)

func notImplementedConsumeFn(*kafkalib.Message) error {
	panic("consume method not implemented")
}

func New(name string) Consumer {
	return &consumerImpl{
		name: name,
		fn:   notImplementedConsumeFn,
	}
}

func NewFn(name string, fn ConsumeFn) Consumer {
	return &consumerImpl{name, fn}
}

func (c *consumerImpl) Name() string {
	return c.name
}

func (c *consumerImpl) Consume(msg *Message) error {
	return c.fn(msg)
}

func (f InterceptorFunc) Intercept(msg *Message, next ConsumeFn) error {
	return f(msg, next)
}
