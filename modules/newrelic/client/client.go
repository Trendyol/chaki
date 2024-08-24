package client_newrelic

import (
	"net/http"

	"github.com/Trendyol/chaki/modules/client/common"
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

func WithClient() common.RoundTripperWrapper {
	return func(rt http.RoundTripper) http.RoundTripper {
		return newRoundTripper(rt)
	}
}
