package main

import (
	"context"

	"github.com/Trendyol/chaki"
	"github.com/Trendyol/chaki/logger"
	"github.com/Trendyol/chaki/modules/client"
	"github.com/Trendyol/chaki/modules/server"
	"github.com/Trendyol/chaki/modules/server/controller"
	"github.com/Trendyol/chaki/modules/server/route"
)

func main() {
	app := chaki.New()

	app.Use(
		client.Module(),
		server.Module(),
	)

	app.Provide(
		newClient,
		NewController,
	)

	if err := app.Start(); err != nil {
		logger.Fatal(err)
	}
}

type serverController struct {
	*controller.Base
	cl *exampleClient
}

func NewController(cl *exampleClient) controller.Controller {
	return &serverController{
		Base: controller.New("server-controller").SetPrefix("/hello"),
		cl:   cl,
	}
}

func (ct *serverController) hello(ctx context.Context, _ any) (string, error) {
	logger.From(ctx).Info("hello")
	resp, err := ct.cl.SendHello(ctx)
	if err != nil {
		return "", err
	}
	return resp, nil
}

func (ct *serverController) greetWithBody(ctx context.Context, req GreetWithBodyRequest) (string, error) {
	return ct.cl.sendGreetWithBody(ctx, req)
}

func (ct *serverController) greetWithQuery(ctx context.Context, req GreetWithQueryRequest) (string, error) {
	return ct.cl.sendGreetWithQuery(ctx, req)
}

func (ct *serverController) greetWithParam(ctx context.Context, req GreetWithParamRequest) (string, error) {
	return ct.cl.sendGreetWithParam(ctx, req)
}

func (ct *serverController) Routes() []route.Route {
	return []route.Route{
		route.Get("", ct.hello),
		route.Get("/query", ct.greetWithQuery).Name("Greet with query"),
		route.Get("/param/:text", ct.greetWithParam).Name("Greet with param"),
		route.Post("/body", ct.greetWithBody).Name("Greet with body"),
	}
}
