package main

import (
	"context"

	"github.com/Trendyol/chaki/modules/server/controller"
	"github.com/Trendyol/chaki/modules/server/response"
	"github.com/Trendyol/chaki/modules/server/route"
)

type CustomRoute struct {
	*controller.Base
	cl *CustomClient
}

func NewCustomController(cl *CustomClient) controller.Controller {
	return &CustomRoute{
		Base: controller.New("client-controller").SetPrefix("/"),
		cl:   cl,
	}
}

func (ct *CustomRoute) Routes() []route.Route {
	return []route.Route{
		route.Post("/:category/", ct.SuccessfulEndpoint),
		route.Post("/:category/error", ct.GetError).Desc("This route has an error from the server itself."),
		route.Post("/:category/error-with-custom-decoder", ct.GetErrNotFound).Desc("This route has error from the custom err decoder."),
		route.Post("/:category/error-with-fallback", ct.GetErrorWithFallback).Desc("This route has a fallback function"),
	}
}

func (ct *CustomRoute) SuccessfulEndpoint(ctx context.Context, req UltimateRequest) (*response.Response[string], error) {
	return ct.cl.SuccessfulEndpoint(ctx, &req)
}

func (ct *CustomRoute) GetError(ctx context.Context, req UltimateRequest) (*response.Response[string], error) {
	return ct.cl.GetError(ctx, &req)
}

func (ct *CustomRoute) GetErrNotFound(ctx context.Context, req UltimateRequest) (*response.Response[string], error) {
	return ct.cl.GetNotFoundErr(ctx, &req)
}

func (ct *CustomRoute) GetErrorWithFallback(ctx context.Context, req UltimateRequest) (*response.Response[string], error) {
	return ct.cl.GetErrorWithFallback(ctx, &req)
}
