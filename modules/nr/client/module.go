package nrclient

import (
	"net/http"

	"github.com/Trendyol/chaki/module"
	clientcommon "github.com/Trendyol/chaki/modules/client/common"
	"github.com/newrelic/go-agent/v3/newrelic"
)

const ModuleName = "chaki-nr-client-module"

func Module() *module.Module {
	return module.New(ModuleName).Provide(newRoundTripper)
}

func newRoundTripper() clientcommon.RoundTripperWrapper {
	return func(rt http.RoundTripper) http.RoundTripper {
		return newrelic.NewRoundTripper(rt)
	}
}
