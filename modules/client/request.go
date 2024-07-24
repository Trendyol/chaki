package client

import (
	"errors"
	"github.com/go-resty/resty/v2"
	"time"
)

var unsupportedMethod = errors.New("unsupported method by the chaki client")

type Request struct {
	*resty.Request

	name string

	// Circuit
	f       circuitFunc
	errF    CircuitErrorFunc
	circuit *circuit
	cmd     string

	// Retry
	*retryConfig
}

type retryConfig struct {
	Retry        int           `json:"retryCount"`
	InitialDelay time.Duration `json:"initialDelay"`
	MaxDelay     time.Duration `json:"maxDelay"`
	Multiplier   float64       `json:"multiplier"`
	DelayType    string        `json:"delayType"`
}

func (r *Request) SetParams(params map[string]string) *Request {
	r.Request.SetQueryParams(params)
	return r
}

func (r *Request) SetParam(k, v string) *Request {
	r.Request.SetQueryParam(k, v)
	return r
}

func (r *Request) SetHeaders(headers map[string]string) *Request {
	r.Request.SetHeaders(headers)
	return r
}

func (r *Request) SetHeader(k, v string) *Request {
	r.Request.SetHeader(k, v)
	return r
}

func (r *Request) SetBody(body any) *Request {
	r.Request.SetBody(body)
	return r
}

func (r *Request) SetResponse(resp any) *Request {
	r.Request.SetResult(resp)
	return r
}

func (r *Request) WithFallback(ef CircuitErrorFunc) *Request {
	r.errF = ef
	return r
}

func (r *Request) WithCommandName(cmd string) *Request {
	r.cmd = cmd
	return r
}

func (r *Request) Post(url string) error {
	r.f = r.functionResolver(url, "POST")
	return r.process()
}

func (r *Request) Get(url string) error {
	r.f = r.functionResolver(url, "GET")
	return r.process()
}

func (r *Request) Delete(url string) error {
	r.f = r.functionResolver(url, "DELETE")
	return r.process()
}

func (r *Request) Put(url string) error {
	r.f = r.functionResolver(url, "PUT")
	return r.process()
}

func (r *Request) Patch(url string) error {
	r.f = r.functionResolver(url, "PATCH")
	return r.process()
}

func (r *Request) process() error {
	err := r.send()
	delay := r.InitialDelay
	for i := 0; i < r.Retry && err != nil; i++ {

		time.Sleep(delay)
		if r.DelayType == "incremental" {
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
	if r.circuit != nil {
		return r.circuit.do(r.cmd, r.f, r.errF)
	}

	return r.f()
}

func (r *Request) functionResolver(url, methodName string) circuitFunc {
	return func() error {
		var err error

		switch methodName {
		case "GET":
			_, err = r.Request.Get(url)
		case "POST":
			_, err = r.Request.Post(url)
		case "PUT":
			_, err = r.Request.Put(url)
		case "PATCH":
			_, err = r.Request.Patch(url)
		case "DELETE":
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
