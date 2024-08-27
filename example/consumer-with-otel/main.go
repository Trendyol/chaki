package main

import (
	"github.com/Trendyol/chaki"
	"github.com/Trendyol/chaki/logger"
	"github.com/Trendyol/chaki/modules/kafka"
	"github.com/Trendyol/chaki/modules/kafka/consumer"
	"github.com/Trendyol/chaki/modules/otel"
	otelkafka "github.com/Trendyol/chaki/modules/otel/kafka"
)

func main() {
	app := chaki.New()

	app.WithOption(
		chaki.WithConfigPath("config.yaml"),
	)

	app.Use(
		otel.Module(
			otelkafka.WithKafka(),
		),
		kafka.Module(),
	)

	app.Provide(
		newConsumer,
	)

	logger.Fatal(app.Start())

}

func newConsumer() consumer.Consumer {
	return consumer.NewFn("example", func(msg *consumer.Message) error {
		// otel span defined in context
		// logger can log spanId and traceId
		logger.From(msg.Context).Info("message here")
		return nil
	})
}
