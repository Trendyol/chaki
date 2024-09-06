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
			ctx := r.Context()
			h := map[string]string{}
			if v := ctxvaluer.CorrelationID.Get(ctx, ""); v != "" {
				h[ctxvaluer.CorrelationIDKey] = v
			}

			if v := ctxvaluer.ExecutorUser.Get(ctx, ""); v != "" {
				h[ctxvaluer.ExecutorUserKey] = v
			}

			if v := ctxvaluer.AgentName.Get(ctx, ""); v != "" {
				h[ctxvaluer.AgentNameKey] = v
			}

			if v := ctxvaluer.Owner.Get(ctx, ""); v != "" {
				h[ctxvaluer.OwnerKey] = v
			}
			r.SetHeaders(h)
			return nil
		})
	}
}
