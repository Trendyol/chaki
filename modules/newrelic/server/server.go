package server

import (
	"github.com/Trendyol/chaki/module"
	nrmodule "github.com/Trendyol/chaki/modules/newrelic"
	"github.com/gofiber/contrib/fibernewrelic"
	"github.com/gofiber/fiber/v2"
	"github.com/newrelic/go-agent/v3/newrelic"
)

func WithServer() nrmodule.Option {
	return nrmodule.WithSubModule(
		module.NewSubModule().Provide(newFiberMiddleware),
	)
}

func newFiberMiddleware(nr *newrelic.Application) fiber.Handler {
	cfg := fibernewrelic.Config{
		Application: nr,
	}

	return fibernewrelic.New(cfg)
}
