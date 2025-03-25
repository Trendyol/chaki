package middlewares

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Trendyol/chaki/util/health"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

// fakeProbe is a mock struct that implements the health.Probe interface.
type fakeProbe struct {
	livenessErr  error
	readinessErr error
}

func (fp *fakeProbe) Liveness(ctx context.Context) error {
	return fp.livenessErr
}

func (fp *fakeProbe) Readiness(ctx context.Context) error {
	return fp.readinessErr
}

// TestHealthCheck verifies the behavior of the HealthCheck middleware using sub-tests.
func TestHealthCheck(t *testing.T) {
	t.Run("it should return success when service is alive and ready", func(t *testing.T) {
		// Both probes return no error.
		probe := &fakeProbe{}
		options := HealthOptions{
			LivenessPath:  "/liveness",
			ReadinessPath: "/readiness",
			Probes:        []health.Probe{probe},
		}

		app := fiber.New()
		app.Use(HealthCheck(options))

		// Test the liveness endpoint.
		req := httptest.NewRequest("GET", "/liveness", nil)
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.NoError(t, resp.Body.Close())

		// Test the readiness endpoint.
		req = httptest.NewRequest("GET", "/readiness", nil)
		resp, err = app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.NoError(t, resp.Body.Close())
	})

	t.Run("it should fail liveness probe when service is not alive", func(t *testing.T) {
		// Make the liveness check return an error.
		probe := &fakeProbe{livenessErr: errors.New("liveness error")}
		options := HealthOptions{
			LivenessPath:  "/liveness",
			ReadinessPath: "/readiness",
			Probes:        []health.Probe{probe},
		}

		app := fiber.New()
		app.Use(HealthCheck(options))

		req := httptest.NewRequest("GET", "/liveness", nil)
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusServiceUnavailable, resp.StatusCode)
		assert.NoError(t, resp.Body.Close())
	})

	t.Run("it should fail liveness probe when service is not ready", func(t *testing.T) {
		// Make the readiness check return an error.
		probe := &fakeProbe{readinessErr: errors.New("readiness error")}
		options := HealthOptions{
			LivenessPath:  "/liveness",
			ReadinessPath: "/readiness",
			Probes:        []health.Probe{probe},
		}

		app := fiber.New()
		app.Use(HealthCheck(options))

		req := httptest.NewRequest("GET", "/readiness", nil)
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusServiceUnavailable, resp.StatusCode)
		assert.NoError(t, resp.Body.Close())
	})
}

// TestRecoverMiddleware tests that the Recover middleware properly handles a panic situation using a sub-test.
func TestRecoverMiddleware(t *testing.T) {
	app := fiber.New()
	app.Use(Recover())

	// An endpoint that will cause a panic.
	app.Get("/panic", func(c *fiber.Ctx) error {
		panic("something bad happened")
	})

	req := httptest.NewRequest("GET", "/panic", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	// The panic should be caught by Recover, resulting in a 500 status code.
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	assert.NoError(t, resp.Body.Close())
}
