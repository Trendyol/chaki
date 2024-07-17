package client

import (
	"github.com/Trendyol/chaki/as"
	"github.com/Trendyol/chaki/module"
	"github.com/Trendyol/chaki/modules/client/common"
	"github.com/go-resty/resty/v2"
	"net/http"
)

var asDriverWrapper = as.Struct[DriverWrapper]("driverupdater")
var asRoundTripperWrapper = as.Struct[common.RoundTripperWrapper]("roundtripperwrapper")

func Module() *module.Module {
	m := module.New("rest")

	m.Provide(
		NewFactory,
		asDriverWrapper.Grouper(),
		asRoundTripperWrapper.Grouper(),
		buildRoundTripperWrapper,
		withCtxBinder,
	)

	m.AddProvideHook(
		module.ProvideHook{
			Match: asDriverWrapper.Match,
			Wrap:  asDriverWrapper.Value,
		},
		module.ProvideHook{
			Match: asRoundTripperWrapper.Match,
			Wrap:  asRoundTripperWrapper.Value,
		},
	)

	return m
}

func buildRoundTripperWrapper(wrappers []common.RoundTripperWrapper) DriverWrapper {
	t := http.DefaultTransport
	for _, wrapper := range wrappers {
		t = wrapper(t)
	}

	return func(c *resty.Client) *resty.Client {
		return c.SetTransport(t)
	}
}
