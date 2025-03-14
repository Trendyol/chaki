package route

import (
	"context"
	"net/http"

	"github.com/Trendyol/chaki/internal/typlect"
	"github.com/Trendyol/chaki/modules/server/response"
	"github.com/Trendyol/chaki/modules/server/validation"
	"github.com/gofiber/fiber/v2"
)

type NoParam = typlect.NoParam

type HandlerFunc[Req, Res any] func(context.Context, Req) (Res, error)

type Route interface {
	Meta() Meta
	Desc(s string) Route
	Name(s string) Route
	AddMiddlewares(mws ...fiber.Handler) Route
}

type route[Req, Res any] struct {
	meta Meta
}

func New[Req, Res any](method, path string, h HandlerFunc[Req, Res], ds ...int) Route {
	return newRoute[Req, Res](method, path, build(h, ds...))
}

func newRoute[Req, Res any](method, path string, h fiber.Handler) Route {
	return &route[Req, Res]{
		meta: Meta{
			Method: method,
			Path:   path,
			Func:   h,
			Req:    typlect.GetType[Req](),
			Res:    typlect.GetType[Res](),
		},
	}
}

func (h *route[Req, Res]) AddMiddlewares(mws ...fiber.Handler) Route {
	h.meta.Middlewares = append(h.meta.Middlewares, mws...)
	return h
}

func (h *route[Req, Res]) Meta() Meta {
	return h.meta
}

func (h *route[Req, Res]) Desc(d string) Route {
	h.meta.Desc = d
	return h
}

func (h *route[Req, Res]) Name(d string) Route {
	h.meta.Name = d
	return h
}

func build[Req, Res any](f HandlerFunc[Req, Res], defaultStatus ...int) fiber.Handler {
	ds := http.StatusOK
	if len(defaultStatus) == 1 {
		ds = defaultStatus[0]
	}

	hasInput := typlect.GetType[Req]() != typlect.TypeNoParam

	return func(c *fiber.Ctx) error {
		var req Req

		hasQuery := len(c.Queries()) > 0
		hasParam := len(c.AllParams()) > 0
		hasContentType := len(c.Request().Header.Peek("content-type")) > 0
		hasContentLength := c.Request().Header.ContentLength() > 0
		hasHeaders := c.Request().Header.Len() > 0
		hasBody := isMethodWithBody(c.Method()) && (hasContentLength || hasContentType || len(c.Body()) > 0)

		if hasInput {
			if hasParam {
				if err := c.ParamsParser(&req); err != nil {
					return err
				}
			}

			if hasQuery {
				if err := c.QueryParser(&req); err != nil {
					return err
				}
			}

			if hasBody && isMethodWithBody(c.Method()) {
				if err := c.BodyParser(&req); err != nil {
					return err
				}
			}

			if hasHeaders {
				if err := c.ReqHeaderParser(&req); err != nil {
					return err
				}
			}

			if err := validation.Validate(req); err != nil {
				return err
			}
		}

		res, err := f(c.UserContext(), req)
		if err != nil {
			return err
		}

		var (
			resp   any
			status = ds
		)

		if rp, ok := any(res).(response.Responser); ok {
			resp = rp.ToResponse()
		} else {
			resp = res
		}

		if st, ok := any(res).(response.Stature); ok {
			status = st.Status()
		}

		return c.Status(status).JSON(resp)
	}
}

func isMethodWithBody(method string) bool {
	switch method {
	case fiber.MethodPost, fiber.MethodPut, fiber.MethodPatch, fiber.MethodDelete:
		return true
	default:
		return false
	}
}
