package middlewares

import (
	"github.com/Trendyol/chaki/logger"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/healthcheck"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"go.uber.org/zap"
	"time"
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

func HealthCheck(liveness, readiness string) fiber.Handler {
	probe := func(c *fiber.Ctx) bool {
		return true
	}

	return healthcheck.New(healthcheck.Config{
		LivenessEndpoint: liveness,
		LivenessProbe:    probe,

		ReadinessEndpoint: readiness,
		ReadinessProbe:    probe,
	})
}

func Recover() fiber.Handler {
	return recover.New()
}
