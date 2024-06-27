package handler

import (
	"context"
	"github.com/Trendyol/chaki/internal/typlect"
	"github.com/Trendyol/chaki/modules/server/response"
	"github.com/Trendyol/chaki/modules/server/validation"
	"github.com/gofiber/fiber/v2"
	"net/http"
)

type NoParam = typlect.NoParam

type HandlerFunc[Req, Res any] func(context.Context, Req) (Res, error)

type Handler interface {
	Meta() Meta
	Desc(s string) Handler
	Name(s string) Handler
}

type handler[Req, Res any] struct {
	meta Meta
}

func New[Req, Res any](method, path string, h HandlerFunc[Req, Res]) Handler {
	return newHandler[Req, Res](method, path, build(h))
}

func newHandler[Req, Res any](method, path string, h fiber.Handler) Handler {
	return &handler[Req, Res]{
		meta: Meta{
			Method: method,
			Path:   path,
			Func:   h,
			Req:    typlect.GetType[Req](),
			Res:    typlect.GetType[Res](),
		},
	}
}

func (h *handler[Req, Res]) AddMiddlewares(mws ...fiber.Handler) Handler {
	h.meta.Middlewares = append(h.meta.Middlewares, mws...)
	return h
}

func (h *handler[Req, Res]) Meta() Meta {
	return h.meta
}

func (h *handler[Req, Res]) Desc(d string) Handler {
	h.meta.Desc = d
	return h
}

func (h *handler[Req, Res]) Name(d string) Handler {
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

			if hasContentType {
				if err := c.BodyParser(&req); err != nil {
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

		if rp, ok := any(res).(responser); ok {
			resp = rp.ToResponse()
		} else {
			resp = response.Success(res)
		}

		if st, ok := any(res).(statuser); ok {
			status = st.Status()
		}

		return c.Status(status).JSON(resp)
	}
}

type responser interface {
	ToResponse() any
}

type statuser interface {
	Status() int
}
