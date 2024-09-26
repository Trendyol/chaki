package main

import (
	"context"
	"fmt"

	"github.com/go-resty/resty/v2"
)

func customErrorDecoder(_ context.Context, res *resty.Response) error {
	if res.StatusCode() == 404 {
		return fmt.Errorf("not found")
	}
	return nil
}
