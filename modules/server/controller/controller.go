package controller

import (
	"github.com/Trendyol/chaki/modules/server/handler"
	"github.com/gofiber/fiber/v2"
)

type Controller interface {
	Handlers() []handler.Handler
	Prefix() string
	Middlewares() []fiber.Handler
	Name() string
}

type Base struct {
	prefix string
	name   string
	mws    []fiber.Handler
}

func New(name string) *Base {
	return &Base{name: name}
}

func (b *Base) Name() string {
	return b.name
}

func (b *Base) Prefix() string {
	return b.prefix
}

func (b *Base) Middlewares() []fiber.Handler {
	return b.mws
}

func (b *Base) Handlers() []handler.Handler {
	panic("handler method not implemented for controller " + b.name)
}

func (b *Base) SetPrefix(p string) *Base {
	b.prefix = p
	return b
}

func (b *Base) AddMiddleware(mw fiber.Handler) *Base {
	b.mws = append(b.mws, mw)
	return b
}
