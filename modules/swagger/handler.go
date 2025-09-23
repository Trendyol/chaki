package swagger

import (
	"net/http"
	"strings"

	"github.com/Trendyol/chaki/config"
	"github.com/Trendyol/chaki/modules/server/common"
	"github.com/Trendyol/chaki/modules/swagger/files"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/fiber/v2/middleware/redirect"
)

func fiberWrapper(docs Docs, cfg *config.Config) common.FiberAppWrapper {
	return func(a *fiber.App) *fiber.App {
		a.Use(
			newHostAccessMiddleware(cfg),
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
			return c.JSON(docs.WithHost(getEffectiveHost(c)))
		}

		if after, ok := strings.CutPrefix(c.Path(), prefix); ok {
			c.Path(after)
			return fsmw(c)
		}

		return c.Next()
	}
}

func newHostAccessMiddleware(cfg *config.Config) fiber.Handler {
	serverCfg := cfg.Of("server")
	if !serverCfg.Exists("swagger") {
		return func(c *fiber.Ctx) error { return c.Next() }
	}

	swcfg := serverCfg.Of("swagger")
	var blocked, allowed []string

	if swcfg.Exists("blockedHostsContains") {
		blocked = swcfg.GetStringSlice("blockedHostsContains")
	}
	if swcfg.Exists("allowedHostsContains") {
		allowed = swcfg.GetStringSlice("allowedHostsContains")
	}

	return func(c *fiber.Ctx) error {
		if !isSwaggerRequest(c.Path()) {
			return c.Next()
		}

		host := getEffectiveHost(c)

		if len(blocked) > 0 && containsAny(host, blocked) {
			return c.SendStatus(fiber.StatusNotFound)
		}

		if len(allowed) > 0 && !containsAny(host, allowed) {
			return c.SendStatus(fiber.StatusNotFound)
		}

		return c.Next()
	}
}

func getEffectiveHost(c *fiber.Ctx) string {
	if xfwd := c.Get("X-Forwarded-Host"); xfwd != "" {
		return xfwd
	}
	return c.Hostname()
}

func isSwaggerRequest(path string) bool {
	if path == "/" || path == "/swagger" || path == "/swagger.json" || path == "/swagger/v1/swagger.json" {
		return true
	}
	return strings.HasPrefix(path, "/swagger/")
}

func containsAny(s string, subs []string) bool {
	for _, sub := range subs {
		if sub != "" && strings.Contains(s, sub) {
			return true
		}
	}
	return false
}
