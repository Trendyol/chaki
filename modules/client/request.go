package client

import (
	"context"
	"errors"
	"github.com/go-resty/resty/v2"
	"time"
)

var unsupportedMethod = errors.New("unsupported method by the chaki client")

const (
	GET = iota
	POST
	PATCH
	PUT
	DELETE
	CUSTOM
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

func (r *Request) Post(url string) error {
	r.f = r.functionResolver(url, POST)
	return r.process()
}

func (r *Request) Get(url string) error {
	r.f = r.functionResolver(url, GET)
	return r.process()
}

func (r *Request) Delete(url string) error {
	r.f = r.functionResolver(url, DELETE)
	return r.process()
}

func (r *Request) Put(url string) error {
	r.f = r.functionResolver(url, PUT)
	return r.process()
}

func (r *Request) Patch(url string) error {
	r.f = r.functionResolver(url, PATCH)
	return r.process()
}

func (r *Request) process() error {
	err := r.send()
	delay := r.Interval
	for i := 0; i < r.Count && err != nil; i++ {

		time.Sleep(delay)
		if r.DelayType == IncrementalDelay {
			delay = time.Duration(float64(delay) * r.Multiplier)
			if delay > r.MaxDelay {
				delay = r.MaxDelay
			}
		}

		select {
		case <-r.Context().Done():
			err = r.Context().Err()
			break
		default:
		}

		err = r.send()
	}

	return err
}

func (r *Request) send() error {
	return r.circuit.do(r.Context(), r.f, r.errF, r.errorFilters...)
}

func (r *Request) functionResolver(url string, method int) CircuitFunc {
	return func(ctx context.Context) error {
		var err error

		switch method {
		case GET:
			_, err = r.Request.Get(url)
		case POST:
			_, err = r.Request.Post(url)
		case PUT:
			_, err = r.Request.Put(url)
		case PATCH:
			_, err = r.Request.Patch(url)
		case DELETE:
			_, err = r.Request.Delete(url)
		default:
			return unsupportedMethod
		}

		if err != nil {
			return err
		}
		return nil
	}
}
