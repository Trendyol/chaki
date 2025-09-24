package main

type UltimateRequest struct {
	Category string `param:"category" validate:"required"`
	Language string `query:"lang" validate:"required,len=2"`
	Message  string `json:"message" validate:"required,min=3,max=100"`
}
