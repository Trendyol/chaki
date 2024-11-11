package client

import (
	"context"
	"errors"
	"math"
	"math/rand"
	"time"

	"github.com/go-resty/resty/v2"
)

var errUnsupportedMethod = errors.New("unsupported method by the chaki client")

const (
	GET = iota
	POST
	PATCH
	PUT
	DELETE
)

type Request struct {
	*resty.Request

	f            CircuitFunc
	errF         CircuitErrorFunc
	errorFilters []CircuitErrorFilter
	circuit      *circuit

	*retryConfig
}

func (r *Request) WithFallback(ef CircuitErrorFunc) *Request {
	r.errF = ef
	return r
}

func (r *Request) WithErrorFilter(f CircuitErrorFilter) *Request {
	r.errorFilters = append(r.errorFilters, f)
	return r
}

func (r *Request) Post(url string) (*resty.Response, error) {
	r.f = r.functionResolver(url, POST)
	return r.process()
}

func (r *Request) Get(url string) (*resty.Response, error) {
	r.f = r.functionResolver(url, GET)
	return r.process()
}

func (r *Request) Delete(url string) (*resty.Response, error) {
	r.f = r.functionResolver(url, DELETE)
	return r.process()
}

func (r *Request) Put(url string) (*resty.Response, error) {
	r.f = r.functionResolver(url, PUT)
	return r.process()
}

func (r *Request) Patch(url string) (*resty.Response, error) {
	r.f = r.functionResolver(url, PATCH)
	return r.process()
}

func (r *Request) process() (*resty.Response, error) {
	resp, err := r.send()
	delay := r.Interval

	for i := 0; i < r.Count && err != nil; i++ {
		if r.DelayType == ExponentialDelay {
			exponentialDelay := delay * time.Duration(math.Pow(2, float64(i)))
			jitter := time.Duration(rand.Float64() * float64(r.Interval))
			delay = exponentialDelay + jitter
			if delay > r.MaxDelay {
				delay = r.MaxDelay
			}
		}

		select {
		case <-r.Context().Done():
			return nil, r.Context().Err()
		case <-time.After(delay):
			resp, err = r.send()
		}
	}

	return resp, err
}

func (r *Request) send() (*resty.Response, error) {
	return r.circuit.do(r.Context(), r.f, r.errF, r.errorFilters...)
}

func (r *Request) functionResolver(url string, method int) CircuitFunc {
	return func(ctx context.Context) (*resty.Response, error) {
		var resp *resty.Response
		var err error

		switch method {
		case GET:
			resp, err = r.Request.Get(url)
		case POST:
			resp, err = r.Request.Post(url)
		case PUT:
			resp, err = r.Request.Put(url)
		case PATCH:
			resp, err = r.Request.Patch(url)
		case DELETE:
			resp, err = r.Request.Delete(url)
		default:
			return nil, errUnsupportedMethod
		}

		if err != nil {
			return nil, err
		}

		return resp, nil
	}
}
