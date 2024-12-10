package client

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/Trendyol/chaki/util/store"
	"github.com/afex/hystrix-go/hystrix"
)

type CircuitRoundTripper struct {
	next     http.RoundTripper
	config   *circuitConfig
	commands *store.Bucket[string, struct{}]
}

func newCircuitRoundTripper(next http.RoundTripper, config *circuitConfig) http.RoundTripper {
	return &CircuitRoundTripper{
		next:     next,
		config:   config,
		commands: store.NewBucket(func(k string) struct{} { return struct{}{} }),
	}
}

func (c *CircuitRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	if c.config == nil || !c.config.Enabled {
		return c.next.RoundTrip(req)
	}

	var circuitName string
	if val := req.Context().Value(circuitCommandKey); val != nil {
		if v, ok := val.(string); ok {
			circuitName = v
		} else {
			return nil, errors.New("circuit command is not valid")
		}
	}

	var fb func(context.Context, error) error
	var fallbackResponse *http.Response
	if val := req.Context().Value(circuitFallbackKey); val != nil {
		if v, ok := val.(func(context.Context, error) (interface{}, error)); ok {
			fb = func(ctx context.Context, err error) error {
				resp, err := v(ctx, err)
				if err != nil {
					return err
				}

				body, contentLength, contentType, err := interfaceToReadCloserWithLength(resp)
				if err != nil {
					return err
				}

				if resp != nil {
					fallbackResponse = &http.Response{
						StatusCode:    200,
						Status:        "200 OK",
						Body:          body,
						Header:        make(http.Header),
						ContentLength: contentLength,
					}

					if contentType != "" {
						fallbackResponse.Header.Set("Content-Type", contentType)
					} else {
						fallbackResponse.Header.Set("Content-Type", "application/json")
					}

					return nil
				}
				return errors.New("could not generate any response from the fallback function" + circuitName)
			}
		} else {
			return nil, errors.New("fallback function is not valid for the circuit: " + circuitName)
		}
	} else {
		fb = defaultCircuitErrorFunc(circuitName)
	}

	var filter func(error) (bool, error)
	if val := req.Context().Value(circuitErrFilterKey); val != nil {
		if v, ok := val.(func(error) (bool, error)); ok {
			filter = v
		} else {
			return nil, errors.New("filter function is not valid for the circuit: " + circuitName)
		}
	}

	var (
		e    error
		ok   bool
		resp *http.Response
	)
	function := func(ctx context.Context) error {
		var err error
		resp, err = c.next.RoundTrip(req)

		if filter != nil {
			if ok, e = filter(err); ok {
				return err
			}
			return nil
		}

		return err
	}

	if !c.commands.Has(circuitName) {
		hystrix.ConfigureCommand(circuitName, c.config.toHystrixConfig())
		c.commands.Set(circuitName, struct{}{})
	}

	hystrixErr := hystrix.DoC(req.Context(), circuitName, function, fb)

	if hystrixErr != nil {
		return nil, hystrixErr
	}

	if e != nil {
		return nil, e
	}

	if fallbackResponse != nil {
		return fallbackResponse, nil
	}

	return resp, nil
}

func interfaceToReadCloserWithLength(data interface{}) (io.ReadCloser, int64, string, error) {
	b, err := json.Marshal(data)
	if err != nil {
		return nil, 0, "", err
	}
	return io.NopCloser(bytes.NewReader(b)), int64(len(b)), "application/json", nil
}
