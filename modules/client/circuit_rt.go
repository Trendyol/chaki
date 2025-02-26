package client

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/Trendyol/chaki/logger"
	"go.uber.org/zap"

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
		return nil, fmt.Errorf("get circuit command on cirucit %s: %w", c.config.Name, err)
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
		return "", fmt.Errorf("circuit %s: command not configured in context", c.config.Name)
	}

	command, ok := val.(string)
	if !ok {
		return "", fmt.Errorf("circuit %s: command must be a string, got %T", c.config.Name, val)
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
		resp *http.Response
		err  error
	)

	execFn := func(ctx context.Context) error {
		resp, err = c.next.RoundTrip(req)

		if err == nil && c.config.shouldTreatStatusCodeAsFailure(resp.StatusCode) {
			respBody := readResponseBody(resp)
			return &GenericClientError{
				c.config.Name,
				resp.StatusCode,
				respBody,
				nil,
			}
		}

		return err
	}

	fbHandler := newOrDefaultFallbackHandler(req.Context())

	if hystrixErr := hystrix.DoC(req.Context(), command, execFn, func(ctx context.Context, errInsideOfFallback error) error {
		err = errInsideOfFallback
		return fbHandler.handle(ctx, err)
	}); hystrixErr != nil {
		return nil, hystrixErr
	}

	if fbHandler.executed {
		logger.From(req.Context()).Warn("fallback executed",
			zap.String("command", command),
			zap.Int("status_code", getStatusCode(resp)),
			zap.String("error_type", getErrorType(err)),
			zap.Error(err))
		return fbHandler.resp, nil
	}

	return resp, err
}

func getErrorType(err error) string {
	var statusErr *GenericClientError

	switch {
	case errors.As(err, &statusErr):
		return "status_code_error"
	case errors.Is(err, hystrix.ErrCircuitOpen):
		return "circuit_open"
	case errors.Is(err, hystrix.ErrTimeout):
		return "timeout"
	case errors.Is(err, hystrix.ErrMaxConcurrency):
		return "max_concurrency"
	default:
		return "other_error"
	}
}

func getStatusCode(resp *http.Response) int {
	if resp == nil {
		return 0
	}
	return resp.StatusCode
}

func readResponseBody(resp *http.Response) []byte {
	if resp == nil || resp.Body == nil {
		return nil
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil
	}
	resp.Body.Close()
	resp.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	return bodyBytes
}
