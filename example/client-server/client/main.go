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
		newExampleClientWithNoRetryAndCircuitBreaker,
		newExampleClientWithRetryAndNoCircuitBreaker,
		newExampleClientWithCircuitBreakerAndNoRetry,
		newExampleClientWithRetryAndCircuitBreaker,
		NewController,
	)

	if err := app.Start(); err != nil {
		logger.Fatal(err)
	}
}

type serverController struct {
	*controller.Base
	exampleClientWithRetryAndCircuitBreaker     *exampleClientWithRetryAndCircuitBreaker
	exampleClientWithNoRetryAndNoCircuitBreaker *exampleClientWithNoRetryAndNoCircuitBreaker
	exampleClientWithCircuitBreakerAndNoRetry   *exampleClientWithCircuitBreakerAndNoRetry
	exampleClientWithRetryAndNoCircuitBreaker   *exampleClientWithRetryAndNoCircuitBreaker
}

func NewController(exampleClientWithRetryAndCircuitBreaker *exampleClientWithRetryAndCircuitBreaker,
	exampleClientWithNoRetryAndNoCircuitBreaker *exampleClientWithNoRetryAndNoCircuitBreaker,
	exampleClientWithCircuitBreakerAndNoRetry *exampleClientWithCircuitBreakerAndNoRetry,
	exampleClientWithRetryAndNoCircuitBreaker *exampleClientWithRetryAndNoCircuitBreaker) controller.Controller {
	return &serverController{
		Base:                                    controller.New("server-controller").SetPrefix("/hello"),
		exampleClientWithRetryAndCircuitBreaker: exampleClientWithRetryAndCircuitBreaker,
		exampleClientWithNoRetryAndNoCircuitBreaker: exampleClientWithNoRetryAndNoCircuitBreaker,
		exampleClientWithCircuitBreakerAndNoRetry:   exampleClientWithCircuitBreakerAndNoRetry,
		exampleClientWithRetryAndNoCircuitBreaker:   exampleClientWithRetryAndNoCircuitBreaker,
	}
}

func (ct *serverController) withRetryAndCircuitBreaker_hello(ctx context.Context, _ any) (string, error) {
	logger.From(ctx).Info("hello")
	resp, err := ct.exampleClientWithRetryAndCircuitBreaker.SendHello(ctx)
	if err != nil {
		return "", err
	}
	return resp, nil
}

func (ct *serverController) withRetryAndCircuitBreaker_greetWithBody(ctx context.Context, req GreetWithBodyRequest) (string, error) {
	return ct.exampleClientWithRetryAndCircuitBreaker.sendGreetWithBody(ctx, req)
}

func (ct *serverController) withRetryAndCircuitBreaker_greetWithQuery(ctx context.Context, req GreetWithQueryRequest) (string, error) {
	return ct.exampleClientWithRetryAndCircuitBreaker.sendGreetWithQuery(ctx, req)
}

func (ct *serverController) withRetryAndCircuitBreaker_greetWithParam(ctx context.Context, req GreetWithParamRequest) (string, error) {
	return ct.exampleClientWithRetryAndCircuitBreaker.sendGreetWithParam(ctx, req)
}

func (ct *serverController) withNoRetryAndNoCircuitBreaker_hello(ctx context.Context, _ any) (string, error) {
	logger.From(ctx).Info("hello")
	resp, err := ct.exampleClientWithNoRetryAndNoCircuitBreaker.SendHello(ctx)
	if err != nil {
		return "", err
	}
	return resp, nil
}

func (ct *serverController) withNoRetryAndNoCircuitBreaker_greetWithBody(ctx context.Context, req GreetWithBodyRequest) (string, error) {
	return ct.exampleClientWithNoRetryAndNoCircuitBreaker.sendGreetWithBody(ctx, req)
}

func (ct *serverController) withNoRetryAndNoCircuitBreaker_greetWithQuery(ctx context.Context, req GreetWithQueryRequest) (string, error) {
	return ct.exampleClientWithNoRetryAndNoCircuitBreaker.sendGreetWithQuery(ctx, req)
}

func (ct *serverController) withNoRetryAndNoCircuitBreaker_greetWithParam(ctx context.Context, req GreetWithParamRequest) (string, error) {
	return ct.exampleClientWithNoRetryAndNoCircuitBreaker.sendGreetWithParam(ctx, req)
}

func (ct *serverController) withCircuitBreakerAndNoRetry_hello(ctx context.Context, _ any) (string, error) {
	logger.From(ctx).Info("hello")
	resp, err := ct.exampleClientWithCircuitBreakerAndNoRetry.SendHello(ctx)
	if err != nil {
		return "", err
	}
	return resp, nil
}

func (ct *serverController) withCircuitBreakerAndNoRetry_greetWithBody(ctx context.Context, req GreetWithBodyRequest) (string, error) {
	return ct.exampleClientWithCircuitBreakerAndNoRetry.sendGreetWithBody(ctx, req)
}

func (ct *serverController) withCircuitBreakerAndNoRetry_greetWithQuery(ctx context.Context, req GreetWithQueryRequest) (string, error) {
	return ct.exampleClientWithCircuitBreakerAndNoRetry.sendGreetWithQuery(ctx, req)
}

func (ct *serverController) withCircuitBreakerAndNoRetry_greetWithParam(ctx context.Context, req GreetWithParamRequest) (string, error) {
	return ct.exampleClientWithCircuitBreakerAndNoRetry.sendGreetWithParam(ctx, req)
}

func (ct *serverController) withNoCircuitBreakerAndRetry_hello(ctx context.Context, _ any) (string, error) {
	logger.From(ctx).Info("hello")
	resp, err := ct.exampleClientWithRetryAndNoCircuitBreaker.SendHello(ctx)
	if err != nil {
		return "", err
	}
	return resp, nil
}

func (ct *serverController) withNoCircuitBreakerAndRetry_greetWithBody(ctx context.Context, req GreetWithBodyRequest) (string, error) {
	return ct.exampleClientWithRetryAndNoCircuitBreaker.sendGreetWithBody(ctx, req)
}

func (ct *serverController) withNoCircuitBreakerAndRetry_greetWithQuery(ctx context.Context, req GreetWithQueryRequest) (string, error) {
	return ct.exampleClientWithRetryAndNoCircuitBreaker.sendGreetWithQuery(ctx, req)
}

func (ct *serverController) withNoCircuitBreakerAndRetry_greetWithParam(ctx context.Context, req GreetWithParamRequest) (string, error) {
	return ct.exampleClientWithRetryAndNoCircuitBreaker.sendGreetWithParam(ctx, req)
}

func (ct *serverController) Routes() []route.Route {
	return []route.Route{
		route.Get("/withRetryAndCircuitBreaker", ct.withRetryAndCircuitBreaker_hello),
		route.Get("/withRetryAndCircuitBreaker/query", ct.withRetryAndCircuitBreaker_greetWithQuery).Name("[Retry And Circuit Breaker] Greet with query"),
		route.Get("/withRetryAndCircuitBreaker/param/:text", ct.withRetryAndCircuitBreaker_greetWithParam).Name("[Retry And Circuit Breaker] Greet with param"),
		route.Post("/withRetryAndCircuitBreaker/body", ct.withRetryAndCircuitBreaker_greetWithBody).Name("[Retry And Circuit Breaker] Greet with body"),

		route.Get("/withNoRetryAndNoCircuitBreaker", ct.withNoRetryAndNoCircuitBreaker_hello),
		route.Get("/withNoRetryAndNoCircuitBreaker/query", ct.withNoRetryAndNoCircuitBreaker_greetWithQuery).Name("[No Retry And No Circuit Breaker] Greet with query"),
		route.Get("/withNoRetryAndNoCircuitBreaker/param/:text", ct.withNoRetryAndNoCircuitBreaker_greetWithParam).Name("[No Retry And No Circuit Breaker] Greet with param"),
		route.Post("/withNoRetryAndNoCircuitBreaker/body", ct.withNoRetryAndNoCircuitBreaker_greetWithBody).Name("[No Retry And No Circuit Breaker] Greet with body"),

		route.Get("/withCircuitBreakerAndNoRetry", ct.withCircuitBreakerAndNoRetry_hello),
		route.Get("/withCircuitBreakerAndNoRetry/query", ct.withCircuitBreakerAndNoRetry_greetWithQuery).Name("[No Retry But With Circuit Breaker] Greet with query"),
		route.Get("/withCircuitBreakerAndNoRetry/param/:text", ct.withCircuitBreakerAndNoRetry_greetWithParam).Name("[No Retry But With Circuit Breaker] Greet with param"),
		route.Post("/withCircuitBreakerAndNoRetry/body", ct.withCircuitBreakerAndNoRetry_greetWithBody).Name("[No Retry But With Circuit Breaker] Greet with body"),

		route.Get("/withNoCircuitBreakerAndRetry", ct.withNoCircuitBreakerAndRetry_hello),
		route.Get("/withCircuitBreakerAndNoRetry/query", ct.withNoCircuitBreakerAndRetry_greetWithQuery).Name("[Retry But No Circuit Breaker] Greet with query"),
		route.Get("/withCircuitBreakerAndNoRetry/param/:text", ct.withNoCircuitBreakerAndRetry_greetWithParam).Name("[Retry But No Circuit Breaker] Greet with param"),
		route.Post("/withCircuitBreakerAndNoRetry/body", ct.withNoCircuitBreakerAndRetry_greetWithBody).Name("[Retry But No Circuit Breaker] Greet with body"),
	}
}
