package client

import (
	"context"
)

// Get to be able to return the response, should pass the address of the response variable.
func (r *Base) Get(ctx context.Context, url string, params map[string]string, res any) error {
	resp, err := r.driver.R().SetContext(ctx).SetQueryParams(params).Get(url)
	if err != nil {
		return err
	}

	if resp.IsSuccess() && resp.Size() > 0 {
		err = r.driver.JSONUnmarshal(resp.Body(), res)
		if err != nil {
			return err
		}
	}

	return nil
}

// Post to be able to return the response, should pass the address of the response variable.
func (r *Base) Post(ctx context.Context, url string, params map[string]string, body any, res any) error {
	resp, err := r.driver.R().SetContext(ctx).SetQueryParams(params).SetBody(body).Post(url)
	if err != nil {
		return err
	}

	if resp.IsSuccess() && resp.Size() > 0 {
		err = r.driver.JSONUnmarshal(resp.Body(), res)
		if err != nil {
			return err
		}
	}

	return nil
}

// Put to be able to return the response, should pass the address of the response variable.
func (r *Base) Put(ctx context.Context, url string, params map[string]string, body any, res any) error {
	resp, err := r.driver.R().SetContext(ctx).SetQueryParams(params).SetBody(body).Put(url)
	if err != nil {
		return err
	}

	if resp.IsSuccess() && resp.Size() > 0 {
		err = r.driver.JSONUnmarshal(resp.Body(), res)
		if err != nil {
			return err
		}
	}

	return nil
}

// Patch to be able to return the response, should pass the address of the response variable.
func (r *Base) Patch(ctx context.Context, url string, params map[string]string, body any, res any) error {
	resp, err := r.driver.R().SetContext(ctx).SetQueryParams(params).SetBody(body).Patch(url)
	if err != nil {
		return err
	}

	if resp.IsSuccess() && resp.Size() > 0 {
		err = r.driver.JSONUnmarshal(resp.Body(), res)
		if err != nil {
			return err
		}
	}

	return nil
}
