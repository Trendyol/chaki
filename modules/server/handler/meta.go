package handler

import (
	"github.com/gofiber/fiber/v2"
	"reflect"
)

type Meta struct {
	Method        string
	Path          string
	Name          string
	Desc          string
	Func          fiber.Handler
	Req           reflect.Type
	Res           reflect.Type
	Middlewares   []fiber.Handler
	DefaultStatus int
}
