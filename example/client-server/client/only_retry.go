package main

import (
	"context"
	"fmt"
	"github.com/Trendyol/chaki/modules/client"
	"github.com/go-resty/resty/v2"
	"net/http"
)

type exampleClientWithRetryAndNoCircuitBreaker struct {
	*client.Base
}

func newExampleClientWithRetryAndNoCircuitBreaker(f *client.Factory) *exampleClientWithRetryAndNoCircuitBreaker {
	return &exampleClientWithRetryAndNoCircuitBreaker{
		Base: f.Get("example-client-with-retry-and-no-circuit-breaker", client.DefaultErrDecoder),
	}
}

func (cl *exampleClientWithRetryAndNoCircuitBreaker) SendHello(ctx context.Context) (string, error) {
	request := cl.Request(ctx)
	var resp *resty.Response
	var err error
	if resp, err = cl.Send(request, "/hello", http.MethodGet); err != nil {
		return "", err
	}

	return string(resp.Body()), nil
}

func (cl *exampleClientWithRetryAndNoCircuitBreaker) sendGreetWithQuery(ctx context.Context, req GreetWithQueryRequest) (string, error) {
	params := map[string]string{
		"text":        req.Text,
		"repeatTimes": fmt.Sprintf("%d", req.RepeatTimes),
	}
	request := cl.Request(ctx).
		SetQueryParams(params)

	var resp *resty.Response
	var err error
	if resp, err = cl.Send(request, "/hello/query", http.MethodGet); err != nil {
		return "", err
	}

	return string(resp.Body()), nil
}

func (cl *exampleClientWithRetryAndNoCircuitBreaker) sendGreetWithParam(ctx context.Context, req GreetWithParamRequest) (string, error) {
	url := fmt.Sprintf("/hello/param/%s", req.Text)
	params := map[string]string{
		"text":        req.Text,
		"repeatTimes": fmt.Sprintf("%d", req.RepeatTimes),
	}
	request := cl.Request(ctx).
		SetQueryParams(params)

	var resp *resty.Response
	var err error
	if resp, err = cl.Send(request, url, http.MethodGet); err != nil {
		return "", err
	}

	return string(resp.Body()), nil
}

func (cl *exampleClientWithRetryAndNoCircuitBreaker) sendGreetWithBody(ctx context.Context, req GreetWithBodyRequest) (string, error) {
	var resp *resty.Response
	request := cl.Request(ctx).
		SetBody(req)

	var err error
	if resp, err = cl.Send(request, "/hello/body", http.MethodPost); err != nil {
		return "", err
	}

	return string(resp.Body()), nil
}
