package main

import (
	"context"

	"github.com/Trendyol/chaki"
	"github.com/Trendyol/chaki/logger"
	"github.com/Trendyol/chaki/modules/kafka"
	"github.com/Trendyol/chaki/modules/kafka/producer"
)

func main() {
	app := chaki.New()

	app.WithOption(
		chaki.WithConfigPath("config.yaml"),
	)

	app.Use(
		kafka.Module(),
	)

	app.Provide(newInterceptor)

	app.Invoke(produceMsg)
}

func newInterceptor() producer.Interceptor {
	return producer.InterceptorFunc(func(ctx context.Context, msgs []producer.Message, next producer.InterceptNextFunc) error {
		logger.From(ctx).Info("interceptor before")
		err := next(ctx, msgs...)
		logger.From(ctx).Info("interceptor after")
		return err
	})
}

func produceMsg(pf producer.Factory) error {
	return pf.Get("example").Produce(context.Background(), producer.Message{
		Value: []byte("test-msg"),
	})
}
