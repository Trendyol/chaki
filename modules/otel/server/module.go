package otelserver

import (
	"github.com/Trendyol/chaki/module"
	otelcommon "github.com/Trendyol/chaki/modules/otel/common"
	servercommon "github.com/Trendyol/chaki/modules/server/common"
	"github.com/gofiber/contrib/otelfiber/v2"
	"github.com/gofiber/fiber/v2"
	"go.opentelemetry.io/otel/trace"
)

const ModuleName = "chaki-otel-server"

func Module() *module.Module {
	return module.New(ModuleName).Provide(newMiddlewareGroup)
}

func newMiddlewareGroup() servercommon.MiddlewareGroup {
	return servercommon.MiddlewareGroup{
		otelfiber.Middleware(),
		ctxSetterMiddleware(),
	}
}

func ctxSetterMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		span := trace.SpanFromContext(c.UserContext())
		ctx := otelcommon.WithContext(c.UserContext(), span)
		c.SetUserContext(ctx)
		return c.Next()
	}
}
