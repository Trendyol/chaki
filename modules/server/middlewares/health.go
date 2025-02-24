package middlewares

import (
	"github.com/Trendyol/chaki/logger"
	"github.com/Trendyol/chaki/util/health"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/healthcheck"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"go.uber.org/zap"
)

type HealthOptions struct {
	LivenessPath  string
	ReadinessPath string
	Probes        []health.Probe
}

func (ho HealthOptions) livenessProbe(c *fiber.Ctx) bool {
	for _, ch := range ho.Probes {
		if err := ch.Liveness(c.UserContext()); err != nil {
			logger.From(c.UserContext()).Error("liveness check failed", zap.Error(err))
			return false
		}
	}
	return true
}

func (ho HealthOptions) readinessProbe(c *fiber.Ctx) bool {
	for _, ch := range ho.Probes {
		if err := ch.Readiness(c.UserContext()); err != nil {
			logger.From(c.UserContext()).Error("liveness check failed", zap.Error(err))
			return false
		}
	}
	return true
}

func HealthCheck(options HealthOptions) fiber.Handler {
	return healthcheck.New(healthcheck.Config{
		LivenessEndpoint:  options.LivenessPath,
		LivenessProbe:     options.livenessProbe,
		ReadinessEndpoint: options.ReadinessPath,
		ReadinessProbe:    options.readinessProbe,
	})
}

func Recover() fiber.Handler {
	return recover.New()
}
