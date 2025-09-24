package main

import "fmt"

type UltimateRequest struct {
	Category string `param:"category" validate:"required"`
	Language string `query:"lang" validate:"required,len=2"`
	Message  string `json:"message" validate:"required,min=3,max=100"`
}

func (ur *UltimateRequest) ToResponse() string {
	return fmt.Sprintf("Category: %s, Language: %s, Message: %s", ur.Category, ur.Language, ur.Message)
}
