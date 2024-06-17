package main

import (
	"github.com/Trendyol/chaki"
	"github.com/Trendyol/chaki/logger"
	"github.com/Trendyol/chaki/modules/kafka"
	"github.com/Trendyol/chaki/modules/kafka/consumer"
)

func main() {
	app := chaki.New()

	app.WithOption(
		chaki.WithConfigPath("config.yaml"),
	)

	app.Use(
		kafka.Module(),
	)

	app.Provide(
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
	logger.From(msg.Context).Info("message arrived")
	return nil
}

func NewExampleConsumer() consumer.Consumer {
	return consumer.NewFn("example1", func(m *consumer.Message) error {
		logger.From(m.Context).Info("message arrived")
		return nil
	})
}
