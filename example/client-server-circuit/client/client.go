package main

import (
	"context"
	"fmt"
	"strconv"

	"github.com/Trendyol/chaki/modules/client"
	"github.com/Trendyol/chaki/modules/server/response"
)

type exampleClient struct {
	*client.Base
}

func newClient(f *client.Factory) *exampleClient {
	return &exampleClient{
		Base: f.Get("example-client", client.WithDriverWrappers(HeaderWrapper())),
	}
}

func (cl *exampleClient) SendHello(ctx context.Context) (string, error) {
	resp := &response.Response[string]{}
	if _, err := cl.Request(ctx).SetResult(resp).Get("/hello"); err != nil {
		return "", err
	}

	return resp.Data, nil
}

func (cl *exampleClient) sendGreetWithQuery(ctx context.Context, req GreetWithQueryRequest) (string, error) {
	resp := ""

	params := map[string]string{
		"text":        req.Text,
		"repeatTimes": fmt.Sprintf("%d", req.RepeatTimes),
	}

	r, err := cl.Request(ctx).
		SetResult(resp).
		SetQueryParams(params).
		SetHeaders(params).
		Get("/hello/query")
	_ = r
	if err != nil {
		return "", err
	}
	return resp, nil
}

func (cl *exampleClient) sendGreetWithParam(ctx context.Context, req GreetWithParamRequest) (string, error) {
	resp := &response.Response[string]{}

	ctx = client.SetFallbackFunc(ctx, func(ctx context.Context, err error) (interface{}, error) {
		return response.Success("custom fallback response"), nil
	})

	if _, err := cl.RequestWithCommand(ctx, "param").
		SetResult(resp).
		SetQueryParam("repeatTimes", strconv.Itoa(req.RepeatTimes)).
		SetPathParam("text", req.Text).
		Get("/hello/param/{text}"); err != nil {
		return "", err
	}

	return resp.Data, nil
}

func (cl *exampleClient) sendGreetWithBody(ctx context.Context, req GreetWithBodyRequest) (string, error) {
	resp := &response.Response[string]{}

	if _, err := cl.Request(ctx).
		SetResult(resp).
		SetBody(req).
		Post("/hello/body"); err != nil {
		return "", err
	}

	return resp.Data, nil
}
