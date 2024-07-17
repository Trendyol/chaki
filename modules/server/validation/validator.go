package validation

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type (
	FieldError struct {
		FailedField string      `json:"failedField"`
		Tag         string      `json:"tag"`
		Value       interface{} `json:"value"`
		Message     string      `json:"message"`
	}

	Error struct {
		msg  string
		errs []FieldError
	}
)

var (
	validate           *validator.Validate
	validationMessages = map[string]func(string, any) string{
		"required": func(s string, a any) string {
			return fmt.Sprintf("field %s is required", s)
		},
	}
)

var structFieldTags = []string{"json", "param", "query"}

func getStructFieldName(fld reflect.StructField) string {
	for _, st := range structFieldTags {
		name := strings.SplitN(fld.Tag.Get(st), ",", 2)[0]

		if name == "" {
			continue
		}

		if name == "-" {
			return ""
		}

		return name
	}

	return fld.Name
}

func Init(rules []Rule) {
	validate = validator.New()

	validate.RegisterTagNameFunc(getStructFieldName)

	for _, r := range rules {
		if err := validate.RegisterValidation(r.Name(), r.Apply); err != nil {
			panic(err.Error())
		}
		validationMessages[r.Name()] = r.Message
	}

}

func Validate(s any) error {
	if validate == nil {
		panic("validator not initialized")
	}

	errs := validate.Struct(s)
	if errs != nil {
		var fieldErrors []FieldError

		var validatorErrors validator.ValidationErrors
		errors.As(errs, &validatorErrors)

		for _, err := range validatorErrors {
			fieldErrors = append(fieldErrors, FieldError{
				FailedField: err.Field(),
				Tag:         err.Tag(),
				Value:       err.Value(),
				Message:     getTagMessage(err),
			})
		}

		return &Error{
			msg:  "validation error",
			errs: fieldErrors,
		}

	}

	return nil
}

func (ve *Error) Unwrap() error {
	return fiber.ErrBadRequest
}

func (ve *Error) Error() string {
	return ve.msg
}

func (ve *Error) FieldErrors() []FieldError {
	return ve.errs
}

func (fe FieldError) Error() string {
	return fe.Message
}

func getTagMessage(err validator.FieldError) string {
	if mr, ok := validationMessages[err.Tag()]; ok {
		return mr(err.Field(), err.Value())
	}

	return err.Error()
}
