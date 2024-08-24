package server_newrelic

import (
	"github.com/gofiber/contrib/fibernewrelic"
	"github.com/gofiber/fiber/v2"
	"github.com/newrelic/go-agent/v3/newrelic"
)

func WithServer() any {
	return newFiberMiddleware
}

func newFiberMiddleware(nr *newrelic.Application) fiber.Handler {
	cfg := fibernewrelic.Config{
		Application: nr,
	}

	return fibernewrelic.New(cfg)
}
