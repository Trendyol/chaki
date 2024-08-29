package main

import (
	"context"
	"fmt"

	"github.com/Trendyol/chaki"
	"github.com/Trendyol/chaki/logger"
	"github.com/Trendyol/chaki/modules/kafka"
	"github.com/Trendyol/chaki/modules/kafka/consumer"
)

func main() {
	app := chaki.New()

	app.WithOption(
		chaki.WithConfigPath("example/batch-consumer/config.yaml"),
	)

	app.Use(
		kafka.Module(),
	)

	app.Provide(
		newInterceptor,
		NewExample2Consumer,
	)

	logger.Fatal(app.Start())
}

type example2Consumer struct {
	consumer.BatchConsumer
}

func NewExample2Consumer() consumer.BatchConsumer {
	return &example2Consumer{
		BatchConsumer: consumer.NewBatchConsumer("example2"),
	}
}

func (c *example2Consumer) Consume(messages []*consumer.Message) error {
	logger.From(context.Background()).Info(fmt.Sprintf("Message Count: %d", len(messages)))

	for i, message := range messages {
		logger.From(message.Context).Info(fmt.Sprintf("message arrived. message order: %d , message: %+v", i, message))
	}

	return nil
}

func newInterceptor() consumer.BatchConsumerInterceptor {
	return consumer.BatchConsumerInterceptorFunc(func(messages []*consumer.Message, next consumer.BatchConsumeFn) error {
		logger.From(context.Background()).Info("consumer interceptor before")
		err := next(messages)
		logger.From(context.Background()).Info("consumer interceptor after")
		return err
	})
}
