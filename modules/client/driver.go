package client

import (
	"net/http"

	"github.com/Trendyol/chaki/config"
	"github.com/Trendyol/chaki/logger"
	"github.com/Trendyol/chaki/modules/client/common"
	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"
)

type driverBuilder struct {
	cfg        *config.Config
	eh         ErrDecoder
	d          *resty.Client
	updaters   []DriverWrapper
	rtWrappers []common.RoundTripperWrapper
}

func newDriverBuilder(cfg *config.Config) *driverBuilder {
	setDefaults(cfg)

	d := resty.New().
		SetBaseURL(cfg.GetString("baseurl")).
		SetTimeout(cfg.GetDuration("timeout")).

		// Debug mode provides a logging, but it's not in the same format with our logger.
		SetDebug(cfg.GetBool("debug"))
	return &driverBuilder{
		cfg: cfg,
		d:   d,
	}
}

func (b *driverBuilder) AddErrDecoder(eh ErrDecoder) *driverBuilder {
	b.eh = eh
	return b
}

func (b *driverBuilder) AddUpdaters(wrappers ...DriverWrapper) *driverBuilder {
	b.updaters = append(b.updaters, wrappers...)
	return b
}

func (b *driverBuilder) AddRoundTripperWrappers(wrappers ...common.RoundTripperWrapper) *driverBuilder {
	b.rtWrappers = append(b.rtWrappers, wrappers...)
	return b
}

func (b *driverBuilder) SetRetry(retryConfig *retryConfig) *driverBuilder {
	if retryConfig == nil {
		return b
	}

	b.rtWrappers = append(b.rtWrappers, func(rt http.RoundTripper) http.RoundTripper {
		return newRetryRoundTripper(rt, retryConfig)
	})

	return b
}

func (b *driverBuilder) SetCircuit(circuitConfig *circuitConfig) *driverBuilder {
	if circuitConfig == nil {
		return b
	}

	b.rtWrappers = append(b.rtWrappers, func(rt http.RoundTripper) http.RoundTripper {
		return newCircuitRoundTripper(rt, circuitConfig)
	})

	return b
}

func (b *driverBuilder) build() *resty.Client {
	if b.cfg.GetBool("logging") {
		b.useLogging()
	}

	for _, upd := range b.updaters {
		b.d = upd(b.d)
	}

	b.d.SetTransport(b.buildRoundTripper())

	b.d.OnAfterResponse(func(c *resty.Client, r *resty.Response) error {
		return b.eh(r.Request.Context(), r)
	})
	return b.d
}

func (b *driverBuilder) buildRoundTripper() http.RoundTripper {
	rt := b.d.GetClient().Transport
	for _, wr := range b.rtWrappers {
		rt = wr(rt)
	}

	return rt
}

func (b *driverBuilder) useLogging() {
	b.d.OnBeforeRequest(func(c *resty.Client, r *resty.Request) error {
		logger.From(r.Context()).Info(
			"request sent",
			zap.String("method", r.Method),
			zap.String("path", r.URL),
			zap.Any("body", r.Body),
			zap.Any("headers", r.Header),
			zap.Any("query", r.QueryParam),
		)
		return nil
	})

	b.d.OnAfterResponse(func(c *resty.Client, r *resty.Response) error {
		logger.From(r.Request.Context()).Info(
			"response got",
			zap.String("method", r.Request.Method),
			zap.String("path", r.Request.URL),
			zap.Int("status", r.StatusCode()),
			zap.String("body", string(r.Body())),
		)
		return nil
	})
}

func setDefaults(cfg *config.Config) {
	cfg.SetDefault("timeout", "5s")
	cfg.SetDefault("debug", false)
	cfg.SetDefault("logging", false)

	setDefaultCircuitConfigs(cfg)
	setDefaultRetryConfigs(cfg)
}
