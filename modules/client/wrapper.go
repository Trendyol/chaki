package client

import (
	"github.com/Trendyol/chaki/util/wrapper"
	"github.com/go-resty/resty/v2"
	"net/http"
)

type DriverWrapper wrapper.Wrapper[*resty.Client]

func WithTransport(tr http.RoundTripper) DriverWrapper {
	return func(c *resty.Client) *resty.Client {
		return c.SetTransport(tr)
	}
}
