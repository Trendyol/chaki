package main

import (
	"reflect"
	"strings"

	"github.com/Trendyol/chaki/modules/server/validation"
	"github.com/go-playground/validator/v10"
)

func NewHelloTextValidator() validation.Rule {
	return validation.NewRule("textValidator", textValidatorFunc, "the %s field cannot be '%s'")
}

func textValidatorFunc(field validator.FieldLevel) bool {
	switch f := field.Field(); f.Kind() {
	case reflect.String:
		text := f.String()
		if strings.Contains(text, "Good bye") {
			return false
		}
		return true
	default:
		return false
	}
}

func NewQueryTextValidator() validation.Rule {
	return validation.NewRule("queryTextValidator", queryTextValidatorFunc, "the %s field cannot be '%s'")
}

func queryTextValidatorFunc(field validator.FieldLevel) bool {
	switch f := field.Field(); f.Kind() {
	case reflect.String:
		text := f.String()
		if strings.Contains(text, "bye") {
			return false
		}
		return true
	default:
		return false
	}
}
