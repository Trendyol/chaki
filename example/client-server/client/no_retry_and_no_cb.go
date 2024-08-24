package main

import (
	"context"
	"fmt"
	"github.com/Trendyol/chaki/modules/client"
	"github.com/Trendyol/chaki/modules/server/response"
)

type exampleClientWithNoRetryAndNoCircuitBreaker struct {
	*client.Base
}

func newExampleClientWithNoRetryAndCircuitBreaker(f *client.Factory) *exampleClientWithNoRetryAndNoCircuitBreaker {
	return &exampleClientWithNoRetryAndNoCircuitBreaker{
		Base: f.Get("example-client-with-no-retry-and-no-circuit-breaker", client.DefaultErrDecoder),
	}
}

func (cl *exampleClientWithNoRetryAndNoCircuitBreaker) SendHello(ctx context.Context) (string, error) {
	resp := &response.Response[string]{}
	if _, err := cl.Request(ctx).SetResult(resp).Get("/hello"); err != nil {
		return "", err
	}

	return resp.Data, nil
}

func (cl *exampleClientWithNoRetryAndNoCircuitBreaker) sendGreetWithQuery(ctx context.Context, req GreetWithQueryRequest) (string, error) {
	resp := &response.Response[string]{}

	params := map[string]string{
		"text":        req.Text,
		"repeatTimes": fmt.Sprintf("%d", req.RepeatTimes),
	}

	if _, err := cl.Request(ctx).
		SetResult(resp).
		SetQueryParams(params).
		Get("/hello/query"); err != nil {
		return "", err
	}
	return resp.Data, nil
}

func (cl *exampleClientWithNoRetryAndNoCircuitBreaker) sendGreetWithParam(ctx context.Context, req GreetWithParamRequest) (string, error) {
	resp := &response.Response[string]{}

	url := fmt.Sprintf("/hello/param/%s", req.Text)
	params := map[string]string{
		"repeatTimes": fmt.Sprintf("%d", req.RepeatTimes),
	}

	if _, err := cl.Request(ctx).
		SetResult(resp).
		SetQueryParams(params).
		Get(url); err != nil {
		return "", err
	}

	return resp.Data, nil
}

func (cl *exampleClientWithNoRetryAndNoCircuitBreaker) sendGreetWithBody(ctx context.Context, req GreetWithBodyRequest) (string, error) {
	resp := &response.Response[string]{}

	if _, err := cl.Request(ctx).
		SetResult(resp).
		SetBody(req).
		Post("/hello/body"); err != nil {
		return "", err
	}

	return resp.Data, nil
}
