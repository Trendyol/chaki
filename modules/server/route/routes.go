package route

import (
	"github.com/gofiber/fiber/v2"
)

func Get[Req, Res any](path string, h HandlerFunc[Req, Res], ds ...int) Route {
	return New(fiber.MethodGet, path, h, ds...)
}

func Post[Req, Res any](path string, h HandlerFunc[Req, Res], ds ...int) Route {
	return New(fiber.MethodPost, path, h, ds...)
}

func Patch[Req, Res any](path string, h HandlerFunc[Req, Res], ds ...int) Route {
	return New(fiber.MethodPatch, path, h, ds...)
}

func Put[Req, Res any](path string, h HandlerFunc[Req, Res], ds ...int) Route {
	return New(fiber.MethodPut, path, h, ds...)
}

func Delete[Req, Res any](path string, h HandlerFunc[Req, Res], ds ...int) Route {
	return New(fiber.MethodDelete, path, h, ds...)
}
