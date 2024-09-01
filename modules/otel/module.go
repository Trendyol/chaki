package otel

import (
	"github.com/Trendyol/chaki"
	"github.com/Trendyol/chaki/module"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
)

const ModuleName = "chaki-otel"

func Module(opts ...Option) *module.Module {
	o := applyOptions(opts...)

	m := module.New(ModuleName)

	m.Provide(
		chaki.Valuer(o),
		startOtel,
	)

	m.Merge(o.subModules...)

	return m
}

func startOtel(lc fx.Lifecycle, opts *options) (trace.TracerProvider, propagation.TextMapPropagator) {
	if opts.initFunc != nil {
		cf := opts.initFunc()
		lc.Append(fx.StopHook(cf))
	}

	return otel.GetTracerProvider(), otel.GetTextMapPropagator()
}
