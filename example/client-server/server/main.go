package main

import (
	"context"
	"github.com/Trendyol/chaki"
	"github.com/Trendyol/chaki/logger"
	"github.com/Trendyol/chaki/modules/server"
	"github.com/Trendyol/chaki/modules/server/controller"
	"github.com/Trendyol/chaki/modules/server/route"
)

func main() {
	app := chaki.New()

	app.Use(
		server.Module(),
	)

	app.Provide(
		NewController,
	)

	if err := app.Start(); err != nil {
		logger.Fatal(err)
	}

}

type serverController struct {
	*controller.Base
}

func NewController() controller.Controller {
	return &serverController{
		Base: controller.New("server-controller").SetPrefix("/hello"),
	}
}

func (ct *serverController) hello(ctx context.Context, _ any) (string, error) {
	logger.From(ctx).Info("hello from server")
	return "Hi From Server", nil
}

func (ct *serverController) Routes() []route.Route {
	return []route.Route{
		route.Get("/", ct.hello),
		route.Get("/greet", ct.greetHandler).Name("Greet Route"),
		route.Get("/query", ct.greetWithQuery).Name("Greet with query"),
		route.Get("/param/:text", ct.greetWithParam).Name("Greet with param"),
		route.Post("/body", ct.greetWithBody).Name("Greet with body")}
}

func (ct *serverController) greetHandler(_ context.Context, _ struct{}) (string, error) {
	return "Greetings!", nil
}

func (ct *serverController) greetWithBody(_ context.Context, req GreetWithBodyRequest) (string, error) {
	return req.Text, nil
}

func (ct *serverController) greetWithQuery(_ context.Context, req GreetWithQueryRequest) (string, error) {
	return req.Text, nil
}

func (ct *serverController) greetWithParam(_ context.Context, req GreetWithParamRequest) (string, error) {
	return req.Text, nil
}
