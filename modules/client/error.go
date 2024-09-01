package client

import (
	"context"

	"github.com/go-resty/resty/v2"
	"github.com/gofiber/fiber/v2"
)

type ErrDecoder func(context.Context, *resty.Response) error

func DefaultErrDecoder(_ context.Context, res *resty.Response) error {
	if res.StatusCode() > 300 {
		return fiber.NewError(res.StatusCode(), "client error")
	}
	return nil
}
