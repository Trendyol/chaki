package main

import (
	"github.com/Trendyol/chaki/modules/client"
	"github.com/go-resty/resty/v2"
)

type user struct {
	publicUsername string
	publicTag      string
}

func HeaderWrapper() client.DriverWrapper {
	return func(restyClient *resty.Client) *resty.Client {
		return restyClient.OnBeforeRequest(func(c *resty.Client, r *resty.Request) error {
			ctx := r.Context()

			h := map[string]string{}

			if v := ctx.Value("user"); v != nil {
				user := v.(user)
				h["publicUsername"] = user.publicUsername
				h["publicTag"] = user.publicTag
			}

			r.SetHeaders(h)
			return nil
		})
	}
}
