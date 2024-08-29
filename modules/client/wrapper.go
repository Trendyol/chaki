package client

import (
	"net/http"

	"github.com/Trendyol/chaki/modules/common/ctxvaluer"
	"github.com/Trendyol/chaki/util/wrapper"
	"github.com/go-resty/resty/v2"
)

type DriverWrapper wrapper.Wrapper[*resty.Client]

func WithTransport(tr http.RoundTripper) DriverWrapper {
	return func(c *resty.Client) *resty.Client {
		return c.SetTransport(tr)
	}
}

func withCtxBinder() DriverWrapper {
	return func(restyClient *resty.Client) *resty.Client {
		return restyClient.OnBeforeRequest(func(c *resty.Client, r *resty.Request) error {
			context := r.Context()

			h := map[string]string{
				ctxvaluer.CorrelationIDKey: context.Value(ctxvaluer.CorrelationIDKey).(string),
				ctxvaluer.ExecutorUserKey:  context.Value(ctxvaluer.ExecutorUserKey).(string),
				ctxvaluer.AgentNameKey:     context.Value(ctxvaluer.AgentNameKey).(string),
				ctxvaluer.OwnerKey:         context.Value(ctxvaluer.OwnerKey).(string),
			}
			r.SetHeaders(h)

			return nil
		})
	}
}
