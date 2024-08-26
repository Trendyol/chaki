package client

import (
	"net/http"

	"github.com/Trendyol/chaki/module"
	nrmodule "github.com/Trendyol/chaki/modules/newrelic"
	"github.com/newrelic/go-agent/v3/newrelic"
)

type httpRoundTripper struct {
	tr http.RoundTripper
}

func newRoundTripper(tr http.RoundTripper) http.RoundTripper {
	return &httpRoundTripper{tr}
}

func (t *httpRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	txn := newrelic.FromContext(req.Context())
	sgm := newrelic.ExternalSegment{
		StartTime: txn.StartSegmentNow(),
		URL:       req.URL.String(),
	}
	defer sgm.End()

	return t.tr.RoundTrip(req)
}

func WithClient() nrmodule.Option {
	sm := module.NewSubModule().Provide(newRoundTripper)
	return nrmodule.WithSubModule(sm)
}
