package route

import (
	"github.com/gofiber/fiber/v2"
)

func NewFiber[Req, Res any](method, path string, h fiber.Handler) Route {
	return newRoute[Req, Res](method, path, h)
}

func FiberGet[Req, Res any](path string, h fiber.Handler) Route {
	return NewFiber[Req, Res](fiber.MethodGet, path, h)
}

func FiberPost[Req, Res any](path string, h fiber.Handler) Route {
	return NewFiber[Req, Res](fiber.MethodPost, path, h)
}

func FiberDelete[Req, Res any](path string, h fiber.Handler) Route {
	return NewFiber[Req, Res](fiber.MethodDelete, path, h)
}

func FiberPatch[Req, Res any](path string, h fiber.Handler) Route {
	return NewFiber[Req, Res](fiber.MethodPatch, path, h)
}

func FiberPut[Req, Res any](path string, h fiber.Handler) Route {
	return NewFiber[Req, Res](fiber.MethodPut, path, h)
}
