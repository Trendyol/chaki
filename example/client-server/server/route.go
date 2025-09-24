package main

import (
	"context"
	"github.com/Trendyol/chaki/modules/client"
	"github.com/Trendyol/chaki/modules/server/controller"
	"github.com/Trendyol/chaki/modules/server/response"
	"github.com/Trendyol/chaki/modules/server/route"
	"github.com/gofiber/fiber/v2"
)

type CustomRoute struct {
	*controller.Base
}

func NewCustomController() controller.Controller {
	return &CustomRoute{
		Base: controller.New("client-controller").SetPrefix("/"),
	}
}

func (ct *CustomRoute) Routes() []route.Route {
	return []route.Route{
		route.Post("/:category", ct.SuccessfulEndpoint),
		route.Post("/:category/error", ct.GetError),
		route.Post("/:category/not-found", ct.GetNotFound),
	}
}

func (ct *CustomRoute) SuccessfulEndpoint(_ context.Context, req UltimateRequest) (response.Response[string], error) {
	return response.Success(req.ToResponse()), nil
}
func (ct *CustomRoute) GetError(_ context.Context, _ UltimateRequest) (route.NoParam, error) {
	return route.NoParam{}, client.GenericClientError{StatusCode: fiber.StatusServiceUnavailable, ParsedBody: "Some random err, doesn't matter much"}
}

func (ct *CustomRoute) GetNotFound(_ context.Context, _ UltimateRequest) (route.NoParam, error) {
	return route.NoParam{}, client.GenericClientError{StatusCode: fiber.StatusNotFound, ParsedBody: "Some random err, doesn't matter much"}
}
