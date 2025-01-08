package middlewares

import (
	"github.com/Trendyol/chaki/logger"
	"github.com/Trendyol/chaki/util/health"
	"github.com/Trendyol/chaki/util/slc"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/healthcheck"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"go.uber.org/zap"
)

type HealthOptions struct {
	LivenessPath      string
	ReadinessPath     string
	LivenessCheckers  []health.LivenessChecker
	ReadinessCheckers []health.ReadinessChecker
}

func buildProbe(checkers []health.Checker) healthcheck.HealthChecker {
	return func(c *fiber.Ctx) bool {
		for _, ch := range checkers {
			if err := ch.Check(); err != nil {
				logger.From(c.UserContext()).Error(
					"health check failed",
					zap.Error(err),
				)
			}
		}
		return true
	}
}

func HealthCheck(options HealthOptions) fiber.Handler {
	livenessProbe := buildProbe(slc.Map(options.LivenessCheckers, func(ch health.LivenessChecker) health.Checker {
		return ch
	}))

	readinessProbe := buildProbe(slc.Map(options.ReadinessCheckers, func(ch health.ReadinessChecker) health.Checker {
		return ch
	}))

	return healthcheck.New(healthcheck.Config{
		LivenessEndpoint: options.LivenessPath,
		LivenessProbe:    livenessProbe,

		ReadinessEndpoint: options.ReadinessPath,
		ReadinessProbe:    readinessProbe,
	})
}

func Recover() fiber.Handler {
	return recover.New()
}
