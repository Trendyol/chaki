package main

import (
	"context"

	"github.com/Trendyol/chaki"
	"github.com/Trendyol/chaki/logger"
	"github.com/Trendyol/chaki/modules/otel"
	otelserver "github.com/Trendyol/chaki/modules/otel/server"
	"github.com/Trendyol/chaki/modules/server"
	"github.com/Trendyol/chaki/modules/server/controller"
	"github.com/Trendyol/chaki/modules/server/middlewares"
	"github.com/Trendyol/chaki/modules/server/route"
)

func main() {
	app := chaki.New()

	app.WithOption(
		chaki.WithConfigPath("config.yaml"),
	)

	app.Use(
		otel.Module(),
		otelserver.Module(),
		server.Module(),
	)

	app.Provide(
		middlewares.ErrHandler,
		newHelloController,
	)

	if err := app.Start(); err != nil {
		logger.Fatal(err)
	}
}

type HelloController struct {
	*controller.Base
}

func newHelloController() controller.Controller {
	return &HelloController{
		Base: controller.New("Hello Controller").SetPrefix("/hello"),
	}
}

func (c *HelloController) Routes() []route.Route {
	return []route.Route{
		route.Get("/greet", c.greetHandler).Name("Greet Route"),
	}
}

func (c *HelloController) greetHandler(ctx context.Context, _ any) (string, error) {
	logger.From(ctx).Info("trace id and spand id will be logged with message")
	return "Greetings!", nil
}
