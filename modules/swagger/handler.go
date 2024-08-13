package swagger

import (
	"net/http"
	"strings"

	"github.com/Trendyol/chaki/modules/server/common"
	"github.com/Trendyol/chaki/modules/swagger/files"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/fiber/v2/middleware/redirect"
)

func fiberWrapper(docs Docs) common.FiberAppWrapper {
	return func(a *fiber.App) *fiber.App {
		a.Use(
			newRedirectMiddleware(),
			newMiddleware(docs),
		)
		return a
	}
}

func newRedirectMiddleware() fiber.Handler {
	return redirect.New(redirect.Config{
		Rules: map[string]string{
			"/":                        "/swagger/index.html",
			"/swagger":                 "/swagger/index.html",
			"/swagger.json":            "/swagger/docs.json",
			"/swagger/v1/swagger.json": "/swagger/docs.json",
		},
	})
}

func newMiddleware(docs Docs) fiber.Handler {
	fscfg := filesystem.ConfigDefault
	fscfg.Root = http.FS(files.Files)
	fsmw := filesystem.New(fscfg)

	prefix := "/swagger"

	return func(c *fiber.Ctx) error {
		if c.Path() == "/swagger/docs.json" || c.Path() == "/swagger/docs.json/" {
			return c.JSON(docs.WithHost(c.Hostname()))
		}

		if strings.HasPrefix(c.Path(), prefix) {
			c.Path(strings.TrimPrefix(c.Path(), prefix))
			return fsmw(c)
		}

		return c.Next()
	}
}
