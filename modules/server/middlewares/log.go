package middlewares

import (
	"time"

	"github.com/Trendyol/chaki/logger"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func Log() fiber.Handler {
	return func(c *fiber.Ctx) error {
		logger.
			From(c.UserContext()).
			Info(
				"incoming request",
				zap.String("path", c.Path()),
				zap.String("method", c.Method()),
				zap.Any("query", c.Queries()),
				zap.String("body", string(c.Body())),
			)

		startTime := time.Now()
		err := c.Next()

		d := time.Since(startTime)

		logger.
			From(c.UserContext()).
			Info(
				"request executed",
				zap.Duration("executed", d),
			)

		return err
	}
}
