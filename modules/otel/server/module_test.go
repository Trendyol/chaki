package otelserver

import (
	"net/http/httptest"
	"testing"

	"github.com/Trendyol/chaki/modules/common/ctxvaluer"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func Test_newMiddlewareGroup(t *testing.T) {
	t.Run("it should write otel information into context", func(t *testing.T) {
		// Given
		var (
			mg       = newMiddlewareGroup()
			app      = fiber.New()
			req      = httptest.NewRequest(fiber.MethodGet, "/", nil)
			tid, sid = "", ""
		)

		for _, m := range mg {
			app.Use(m)
		}

		app.Use(func(c *fiber.Ctx) error {
			ctx := c.UserContext()
			tid = ctxvaluer.TraceID.Get(ctx)
			sid = ctxvaluer.SpanID.Get(ctx)
			return nil
		})

		// When
		_, err := app.Test(req)

		// Then
		assert.NoError(t, err)
		assert.NotEmpty(t, tid)
		assert.NotEmpty(t, sid)
	})
}
