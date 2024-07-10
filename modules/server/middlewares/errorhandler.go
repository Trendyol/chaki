package middlewares

import (
	"errors"
	"github.com/Trendyol/chaki/logger"
	"github.com/Trendyol/chaki/modules/server/response"
	"github.com/Trendyol/chaki/modules/server/validation"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func ErrHandler(c *fiber.Ctx, err error) error {
	logger.From(c.UserContext()).Error(
		"err handler log",
		zap.Error(err),
	)

	return c.
		Status(getCodeFromErr(err)).
		JSON(response.Error(err.Error(), getValidationErrs(err)...))
}

func getValidationErrs(err error) []validation.FieldError {
	if vErr := new(validation.Error); errors.As(err, &vErr) {
		return vErr.FieldErrors()
	}
	return nil
}

func getCodeFromErr(err error) int {
	if fErr := new(fiber.Error); errors.As(err, &fErr) {
		return fErr.Code
	}

	return fiber.StatusInternalServerError
}
