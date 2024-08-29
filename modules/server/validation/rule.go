package validation

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

type Rule interface {
	Name() string
	Apply(fl validator.FieldLevel) bool
	Message(field string, value any) string
}

type rule struct {
	name     string
	apply    func(fl validator.FieldLevel) bool
	msgTempl string
}

func NewRule(name string, apply func(validator.FieldLevel) bool, msgTempl string) Rule {
	return &rule{
		name:     name,
		apply:    apply,
		msgTempl: msgTempl,
	}
}

func (r *rule) Name() string {
	return r.name
}

func (r *rule) Apply(fl validator.FieldLevel) bool {
	return r.apply(fl)
}

func (r *rule) Message(field string, value any) string {
	return fmt.Sprintf(r.msgTempl, field, value)
}
