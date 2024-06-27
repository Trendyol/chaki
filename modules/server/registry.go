package server

import (
	"github.com/Trendyol/chaki/modules/server/controller"
	"github.com/Trendyol/chaki/modules/server/handler"
	"github.com/Trendyol/chaki/util/slc"
	"github.com/gofiber/fiber/v2"
	"net/url"
)

type registry struct {
	ct       any
	base     string
	name     string
	mws      []fiber.Handler
	handlers []handler.Handler
}

func parseControllers(cts ...controller.Controller) []*registry {
	return slc.Map(cts, newRegistry)
}

func newRegistry(ctr controller.Controller) *registry {
	return &registry{
		ct:       ctr,
		base:     ctr.Prefix(),
		name:     ctr.Name(),
		mws:      ctr.Middlewares(),
		handlers: ctr.Handlers(),
	}
}

func (r *registry) parsePath(path string) string {
	path, err := url.JoinPath(r.base, path)
	if err != nil {
		panic(err.Error())
	}
	return path
}

func (r *registry) toMeta(h handler.Handler) handler.Meta {
	m := h.Meta()
	if m.Name == "" {
		m.Name = r.parsePath(m.Path)
	}
	m.Path = r.parsePath(m.Path)
	return m
}
