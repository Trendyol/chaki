package common

import (
	"github.com/Trendyol/chaki/util/wrapper"
	"github.com/gofiber/fiber/v2"
)

type (
	MiddlewareGroup    = []fiber.Handler
	FiberAppWrapper    = wrapper.Wrapper[*fiber.App]
	FiberConfigWrapper = wrapper.Wrapper[fiber.Config]
)
