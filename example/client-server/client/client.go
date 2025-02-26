package main

import (
	"context"
	"github.com/Trendyol/chaki/modules/client"
	"github.com/Trendyol/chaki/modules/server/response"
	"github.com/go-resty/resty/v2"
)

const (
	errorEndpoint      = "{category}/error"
	notFoundEndpoint   = "{category}/not-found"
	successfulEndpoint = "{category}"
)

type CustomClient struct {
	*client.Base
}

type UltimateRequestBody struct {
	Message string `json:"message"`
}

func NewCustomClient(f *client.Factory) *CustomClient {

	return &CustomClient{
		Base: f.Get("custom-client", client.WithErrDecoder(customErrorDecoder)),
	}
}

func customErrorDecoder(_ context.Context, res *resty.Response) error {
	if res.IsSuccess() {
		return nil
	}

	if res.StatusCode() == 404 {
		return client.GenericClientError{ParsedBody: "not found from custom err decoder", StatusCode: res.StatusCode()}
	}

	return client.GenericClientError{ParsedBody: "generic error se the code :)", StatusCode: res.StatusCode()}
}

func (c *CustomClient) SuccessfulEndpoint(ctx context.Context, req *UltimateRequest) (*response.Response[string], error) {
	resp := &response.Response[string]{}

	if _, err := c.RequestWithCommand(ctx, "commandpostsuccess").
		SetPathParam("category", req.Category).
		SetBody(UltimateRequestBody{Message: req.Message}).
		SetQueryParam("lang", req.Language).
		SetResult(resp).
		Post(successfulEndpoint); err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *CustomClient) GetNotFoundErr(ctx context.Context, req *UltimateRequest) (*response.Response[string], error) {
	resp := &response.Response[string]{}

	if _, err := c.RequestWithCommand(ctx, "commandposterror").
		SetPathParam("category", req.Category).
		SetBody(UltimateRequestBody{Message: req.Message}).
		SetQueryParam("lang", req.Language).
		SetResult(resp).
		Post(notFoundEndpoint); err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *CustomClient) GetError(ctx context.Context, req *UltimateRequest) (*response.Response[string], error) {
	resp := &response.Response[string]{}

	if _, err := c.RequestWithCommand(ctx, "commandposterror").
		SetPathParam("category", req.Category).
		SetBody(UltimateRequestBody{Message: req.Message}).
		SetQueryParam("lang", req.Language).
		SetResult(resp).
		Post(errorEndpoint); err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *CustomClient) GetErrorWithFallback(ctx context.Context, req *UltimateRequest) (*response.Response[string], error) {
	resp := &response.Response[string]{}

	ctx = client.SetFallbackFunc(ctx, func(ctx context.Context, err error) (any, error) {

		// Any fallback mechanism can apply here. The tricky part here is,
		// If you use .SetResult(res) method from resty.Request,
		// You should return the same here. If not, you can handle
		// The fallback response as you wish by using httpRes.Result()
		return &response.Response[string]{
			Data: "this response is from fallback",
		}, nil
	})

	httpRes, err := c.RequestWithCommand(ctx, "commandpostfallback").
		SetPathParam("category", req.Category).
		SetBody(UltimateRequestBody{Message: req.Message}).
		SetQueryParam("lang", "en").
		SetResult(resp).
		Post(errorEndpoint)
	_ = httpRes

	if err != nil {
		return nil, err
	}

	return resp, nil

}
