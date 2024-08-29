package middlewares

import (
	"context"
	"time"

	"github.com/Trendyol/chaki/modules/common/ctxvaluer"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func ContextBinder() fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.SetUserContext(createContext(c))
		return c.Next()
	}
}

func ContextBinderWithTimeout(timeout time.Duration) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx, cancel := context.WithTimeout(createContext(c), timeout)
		defer cancel()

		c.SetUserContext(ctx)

		return c.Next()
	}
}

func createContext(c *fiber.Ctx) context.Context {
	ctx := context.Background()

	return ctxvaluer.CreateBaseTaskContext(ctx, ctxvaluer.CreateParams{
		CorrelationID: c.Get(ctxvaluer.CorrelationIDKey, uuid.NewString()),
		ExecutorUser:  c.Get(ctxvaluer.ExecutorUserKey),
		AgentName:     c.Get(ctxvaluer.AgentNameKey, ""),
		Owner:         c.Get(ctxvaluer.OwnerKey, ""),
	})
}
