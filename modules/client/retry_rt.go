package client

import (
	"math"
	"math/rand"
	"net/http"
	"time"
)

type RetryRoundTripper struct {
	next http.RoundTripper
	cfg  *retryConfig
}

func newRetryRoundTripper(next http.RoundTripper, cfg *retryConfig) http.RoundTripper {
	return &RetryRoundTripper{
		next: next,
		cfg:  cfg,
	}
}

func (r *RetryRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	resp, err := r.next.RoundTrip(req)
	if r.cfg == nil || r.cfg.Count == 0 {
		return resp, err
	}

	delay := r.cfg.Interval
	for i := 0; i < r.cfg.Count && err != nil; i++ {
		if r.cfg.DelayType == ExponentialDelay {
			exponentialDelay := delay * time.Duration(math.Pow(2, float64(i)))

			randFloat := rand.New(rand.NewSource(time.Now().UnixNano())).Float64()
			jitter := time.Duration(randFloat * float64(r.cfg.Interval))
			delay = exponentialDelay + jitter
			if delay > r.cfg.MaxDelay {
				delay = r.cfg.MaxDelay
			}
		}

		select {
		case <-req.Context().Done():
			return nil, req.Context().Err()
		case <-time.After(delay):
			resp, err = r.next.RoundTrip(req)
		}
	}

	return resp, err
}
