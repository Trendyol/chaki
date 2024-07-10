package middlewares

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/healthcheck"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

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
