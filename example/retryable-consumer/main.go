package main

import (
	"errors"
	"fmt"
	"github.com/Trendyol/chaki"
	"github.com/Trendyol/chaki/logger"
	"github.com/Trendyol/chaki/modules/kafka"
	"github.com/Trendyol/chaki/modules/kafka/consumer"
)

func main() {
	app := chaki.New()

	app.WithOption(
		chaki.WithConfigPath("example/retryable-consumer/config.yaml"),
	)

	app.Use(
		kafka.Module(),
	)

	app.Provide(
		newInterceptor,
		NewExample2Consumer,
		NewExampleConsumer,
	)

	logger.Fatal(app.Start())

}

type example2Consumer struct {
	consumer.Consumer
}

func NewExample2Consumer() consumer.Consumer {
	return &example2Consumer{
		Consumer: consumer.New("example2"),
	}
}

func (c *example2Consumer) Consume(msg *consumer.Message) error {
	logger.From(msg.Context).Info(fmt.Sprintf("message arrived. message: %+v", msg))
	return errors.New("error")
}

func NewExampleConsumer() consumer.Consumer {
	return consumer.NewFn("example1", func(m *consumer.Message) error {
		logger.From(m.Context).Info("message arrived")
		return nil
	})
}

func newInterceptor() consumer.Interceptor {
	return consumer.InterceptorFunc(func(msg *consumer.Message, next consumer.ConsumeFn) error {
		logger.From(msg.Context).Info("consumer interceptor before")
		err := next(msg)
		logger.From(msg.Context).Info("consumer interceptor after")
		return err
	})
}
