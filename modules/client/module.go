package client

import (
	"github.com/Trendyol/chaki/as"
	"github.com/Trendyol/chaki/module"
	"github.com/Trendyol/chaki/modules/client/common"
)

var (
	asDriverWrapper       = as.Struct[DriverWrapper]("driverupdater")
	asRoundTripperWrapper = as.Struct[common.RoundTripperWrapper]("roundtripperwrapper")
)

func Module() *module.Module {
	m := module.New("rest")

	m.Provide(
		NewFactory,
		asDriverWrapper.Grouper(),
		asRoundTripperWrapper.Grouper(),
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
