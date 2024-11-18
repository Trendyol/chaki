package client

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/Trendyol/chaki/util/store"
	"github.com/afex/hystrix-go/hystrix"
)

type CircuitRountTripper struct {
	next     http.RoundTripper
	config   *circuitConfig
	commands *store.Bucket[string, struct{}]
}

func newCircuitRoundTripper(next http.RoundTripper, config *circuitConfig) http.RoundTripper {
	return &CircuitRountTripper{
		next:     next,
		config:   config,
		commands: store.NewBucket(func(k string) struct{} { return struct{}{} }),
	}
}

func (c *CircuitRountTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	if c.config == nil || !c.config.Enabled {
		return c.next.RoundTrip(req)
	}

	circuitName := getCircuitName(req)
	var fb func(context.Context, error) error
	if val := req.Context().Value("circuitFallback"); val != nil {
		if v, ok := val.(func(context.Context, error) error); ok {
			fb = v
		} else {
			return nil, errors.New("fallback function is not valid for the circuit: " + circuitName)
		}
	} else {
		fb = defaultCircuitErrorFunc(circuitName)
	}

	var filter func(error) (bool, error)
	if val := req.Context().Value("circuitFilter"); val != nil {
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

	return resp, nil
}

func getCircuitName(req *http.Request) string {
	sb := strings.Builder{}

	sb.WriteString(req.Method)
	sb.WriteString("-")
	sb.WriteString(req.URL.Host)
	sb.WriteString(req.URL.Path)

	return sb.String()
}
