package middlewares

import (
	"context"
	"time"

	"github.com/Trendyol/chaki/config"
	"github.com/Trendyol/chaki/modules/common/ctxvaluer"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func ContextBinder(cfg *config.Config) fiber.Handler {
	customHeaders := cfg.Of("server").GetStringMap("loggingHeaders")
	mapping := ctxvaluer.GetHeaderMapping(customHeaders)

	return func(c *fiber.Ctx) error {
		c.SetUserContext(createContext(c, mapping))
		return c.Next()
	}
}

func ContextBinderWithTimeout(cfg *config.Config, timeout time.Duration) fiber.Handler {
	customHeaders := cfg.Of("server").GetStringMap("loggingHeaders")
	mapping := ctxvaluer.GetHeaderMapping(customHeaders)

	return func(c *fiber.Ctx) error {
		ctx, cancel := context.WithTimeout(createContext(c, mapping), timeout)
		defer cancel()

		c.SetUserContext(ctx)
		return c.Next()
	}
}

func createContext(c *fiber.Ctx, mapping map[string]string) context.Context {
	params := make(ctxvaluer.CreateParams)

	for logKey, headerName := range mapping {
		val := c.Get(headerName)
		if logKey == ctxvaluer.CorrelationIDKey && val == "" {
			val = uuid.NewString()
		}
		params[logKey] = val
	}

	return ctxvaluer.CreateBaseTaskContext(context.Background(), params)
}
