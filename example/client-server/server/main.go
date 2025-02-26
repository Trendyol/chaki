package main

import (
	"context"

	"github.com/Trendyol/chaki"
	"github.com/Trendyol/chaki/modules/otel"
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

		// To add otel module, simply add the following line
		// This requires otel init function and submodules.
		otel.Module(
			otel.WithInitFunc(customOtelInitFunc),
			otelserver.WithServer(),
		),
		swagger.Module(),
	)

	app.Provide(
		NewCustomController,
	)

	_ = app.Start()
}

// You should be setting your propogations, exporters, and other configurations here.
func customOtelInitFunc() otel.CloseFunc {
	return func(ctx context.Context) error {
		return nil
	}
}
