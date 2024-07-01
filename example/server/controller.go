package main

import (
	"context"
	"github.com/Trendyol/chaki/internal/typlect"
	"github.com/Trendyol/chaki/modules/server/controller"
	"github.com/Trendyol/chaki/modules/server/handler"
)

type HelloController struct {
	*controller.Base
	svc *Service
}

func NewHelloController(svc *Service) *HelloController {
	c := controller.New("Hello Controller").SetPrefix("/hello")
	return &HelloController{
		Base: c,
		svc:  svc,
	}
}

func (c *HelloController) Handlers() []handler.Handler {
	return []handler.Handler{
		handler.Get("/greet", c.greetHandler).Name("Greet Handler"),
		handler.Get("/query", c.greetWithQuery).Name("Greet with query"),
		handler.Get("/param/:text", c.greetWithParam).Name("Greet with param"),
		handler.Post("/body", c.greetWithBody).Name("Greet with body"),
	}
}

func (c *HelloController) greetHandler(_ context.Context, _ typlect.NoParam) (string, error) {
	return "Greetings!", nil
}

func (c *HelloController) greetWithBody(_ context.Context, req GreetWithBodyRequest) (string, error) {
	return c.svc.GenerateText(req.Text, req.RepeatTimes)
}

func (c *HelloController) greetWithQuery(_ context.Context, req GreetWithQueryRequest) (string, error) {
	return c.svc.GenerateText(req.Text, req.RepeatTimes)
}

func (c *HelloController) greetWithParam(_ context.Context, req GreetWithParamRequest) (string, error) {
	return c.svc.GenerateText(req.Text, req.RepeatTimes)
}
