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
	if !c.isCircuitEnabled() {
		return c.next.RoundTrip(req)
	}

	command, err := c.getCircuitCommand(req.Context())
	if err != nil {
		return nil, err
	}

	c.ensureCommandConfigured(command)

	return c.executeWithCircuitBreaker(req, command)
}

func (c *CircuitRoundTripper) isCircuitEnabled() bool {
	return c.config != nil && c.config.Enabled
}

func (c *CircuitRoundTripper) getCircuitCommand(ctx context.Context) (string, error) {
	val := ctx.Value(circuitCommandKey)
	if val == nil {
		return "", errors.New("circuit command is not configured")
	}

	command, ok := val.(string)
	if !ok {
		return "", errors.New("circuit command must be a string")
	}

	return command, nil
}

func (c *CircuitRoundTripper) ensureCommandConfigured(command string) {
	if !c.commands.Has(command) {
		hystrix.ConfigureCommand(command, c.config.toHystrixConfig())
		c.commands.Set(command, struct{}{})
	}
}

func (c *CircuitRoundTripper) executeWithCircuitBreaker(req *http.Request, command string) (*http.Response, error) {
	var (
		resp            *http.Response
		fallbackHandler *fallbackHandler
		err             error
	)

	fallbackHandler = newFallbackHandler(req.Context())
	errFilterFunc := getErrorFilterFunc(req.Context())

	execFn := func(ctx context.Context) error {
		resp, err = c.next.RoundTrip(req)

		if modifyErr, filterErr := errFilterFunc(err); modifyErr {
			return filterErr
		}
		return err
	}

	if hystrixErr := hystrix.DoC(req.Context(), command, execFn, fallbackHandler.handle); hystrixErr != nil {
		return nil, hystrixErr
	}

	if err != nil {
		return nil, err
	}

	if fallbackResp := fallbackHandler.resp; fallbackResp != nil {
		return fallbackResp, nil
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
