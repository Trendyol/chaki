package nrserver

import (
	"github.com/Trendyol/chaki/module"
	"github.com/gofiber/contrib/fibernewrelic"
	"github.com/gofiber/fiber/v2"
	"github.com/newrelic/go-agent/v3/newrelic"
)

const ModuleName = "chaki-nr-server-module"

func Module() *module.Module {
	return module.New(ModuleName).Provide(newFiberMiddleware)
}

func newFiberMiddleware(nr *newrelic.Application) fiber.Handler {
	return fibernewrelic.New(fibernewrelic.Config{
		Application: nr,
	})
}
