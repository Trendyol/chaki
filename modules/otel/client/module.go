package otelclient

import (
	"net/http"

	"github.com/Trendyol/chaki/module"
	clientcommon "github.com/Trendyol/chaki/modules/client/common"
	otelmodule "github.com/Trendyol/chaki/modules/otel"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

func WithClient() otelmodule.Option {
	m := module.NewSubModule().Provide(newClientRoundTripperWrapper)
	return otelmodule.WithSubModule(m)
}

func newClientRoundTripperWrapper(tp trace.TracerProvider, tmp propagation.TextMapPropagator) clientcommon.RoundTripperWrapper {
	return func(rt http.RoundTripper) http.RoundTripper {
		return otelhttp.NewTransport(
			rt,
			otelhttp.WithPropagators(tmp),
			otelhttp.WithTracerProvider(tp),
		)
	}
}
