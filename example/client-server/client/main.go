package main

import (
	"context"

	"github.com/Trendyol/chaki"
	"github.com/Trendyol/chaki/modules/client"
	"github.com/Trendyol/chaki/modules/otel"
	otelclient "github.com/Trendyol/chaki/modules/otel/client"
	otelserver "github.com/Trendyol/chaki/modules/otel/server"
	"github.com/Trendyol/chaki/modules/server"
	"github.com/Trendyol/chaki/modules/swagger"
)

func main() {

	app := chaki.New()

	app.WithOption(
		chaki.WithConfigPath("config.yaml"),
	)

	app.Use(
		server.Module(),
		client.Module(),

		// To add otel module, simply add the following line
		// This requires otel init function and submodules.
		otel.Module(
			otel.WithInitFunc(customOtelInitFunc),
			otelserver.WithServer(),
			otelclient.WithClient(),
		),
		swagger.Module(),
	)

	app.Provide(
		NewCustomController,
		NewCustomClient,
	)

	_ = app.Start()
}

// You should be setting your propagations, exporters, and other configurations here.
func customOtelInitFunc() otel.CloseFunc {
	return func(ctx context.Context) error {
		return nil
	}
}
