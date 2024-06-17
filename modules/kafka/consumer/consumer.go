package consumer

import (
	kafkalib "github.com/Trendyol/kafka-konsumer/v2"
)

type (
	Message   = kafkalib.Message
	ConsumeFn = kafkalib.ConsumeFn

	Consumer interface {
		Name() string
		Consume(*Message) error
	}

	consumerImpl struct {
		name string
		fn   ConsumeFn
	}
)

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
