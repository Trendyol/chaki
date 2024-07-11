package response

import "github.com/Trendyol/chaki/modules/server/validation"

type Responser interface {
	ToResponse() any
}

type Statuser interface {
	Status() int
}

type Response[T any] struct {
	Success       bool `json:"success"`
	Data          T    `json:"data,omitempty"`
	Page          *int `json:"page,omitempty"`
	Size          *int `json:"size,omitempty"`
	TotalElements *int `json:"totalElements,omitempty"`
	TotalPages    *int `json:"totalPages,omitempty"`

	// Error fields
	Message          string                  `json:"message,omitempty"`
	ValidationErrors []validation.FieldError `json:"validationErrors,omitempty"`
}

func Success(data any) Response[any] {
	return Response[any]{
		Success: true,
		Data:    data,
	}
}

func Error(msg string, fieldErrs ...validation.FieldError) Response[any] {
	return Response[any]{
		Success:          false,
		Message:          msg,
		ValidationErrors: fieldErrs,
	}
}
