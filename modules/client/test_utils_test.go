// Package client provides HTTP client functionality with circuit breaker, retry, and fallback capabilities.
package client

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"time"

	"github.com/spf13/viper"

	"github.com/Trendyol/chaki/config"
	"github.com/afex/hystrix-go/hystrix"
	"github.com/go-resty/resty/v2"
)

// --- Configuration Utilities ---

// clientTestConfig creates a config with proper structure for client_test.go
// This sets up the config in a nested structure as expected by Factory.Get()
func clientTestConfig() *config.Config {
	cfg := config.NewConfig(viper.New(), nil)
	// Set up values as needed by the client package tests
	// The client package expects a nested structure "client.{clientName}.baseurl"
	cfg.Set("client.testclient.baseurl", "https://example.com")
	cfg.Set("client.testclient.timeout", "1s")
	cfg.Set("client.testclient.debug", false)
	cfg.Set("client.testclient.logging", false)

	// Circuit breaker defaults
	cfg.Set("client.testclient.circuit.enabled", false)
	cfg.Set("client.testclient.circuit.preset", "default")

	// Retry defaults
	cfg.Set("client.testclient.retry.enabled", false)
	cfg.Set("client.testclient.retry.count", 3)
	cfg.Set("client.testclient.retry.waitTime", "100ms")
	cfg.Set("client.testclient.retry.maxWaitTime", "1s")

	return cfg
}

// driverTestConfig creates a config with flat structure for driver_test.go
// This sets up the config in a flat structure as expected by driverBuilder
func driverTestConfig() *config.Config {
	cfg := config.NewConfig(viper.New(), nil)
	// Set up values as needed by the driver builder
	// The driver builder expects values at the root level
	cfg.Set("baseurl", "https://example.com")
	cfg.Set("timeout", "1s")
	cfg.Set("debug", false)
	cfg.Set("logging", false)

	// Circuit breaker defaults
	cfg.Set("circuit.enabled", false)
	cfg.Set("circuit.preset", "default")

	// Retry defaults
	cfg.Set("retry.enabled", false)
	cfg.Set("retry.count", 3)
	cfg.Set("retry.waitTime", "100ms")
	cfg.Set("retry.maxWaitTime", "1s")

	return cfg
}

// --- HTTP Server Mocking ---

// mockServer creates a test HTTP server for integration testing
func mockServer(handler http.HandlerFunc) *httptest.Server {
	return httptest.NewServer(handler)
}

// standardHandler returns a common test handler that returns configurable responses
func standardHandler(statusCode int, body map[string]interface{}, responseDelay time.Duration) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if responseDelay > 0 {
			time.Sleep(responseDelay)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)

		if body != nil {
			json.NewEncoder(w).Encode(body)
		}
	}
}

// --- HTTP Response Utilities ---

// createSuccessResponse generates a successful HTTP response
func createSuccessResponse(req *http.Request, body string) *http.Response {
	if body == "" {
		body = `{"success": true}`
	}

	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(strings.NewReader(body)),
		Header:     http.Header{"Content-Type": []string{"application/json"}},
		Request:    req,
	}
}

// createErrorResponse generates an error HTTP response
// Used in retry_rt_test.go and circuit_rt_test.go
//
//nolint:unused
func createErrorResponse(req *http.Request, statusCode int, body string) *http.Response {
	if body == "" {
		body = `{"error": "Something went wrong"}`
	}

	return &http.Response{
		StatusCode: statusCode,
		Body:       io.NopCloser(strings.NewReader(body)),
		Header:     http.Header{"Content-Type": []string{"application/json"}},
		Request:    req,
	}
}

// --- Mock RoundTripper ---

// MockRoundTripper allows tracking and simulating HTTP requests
type MockRoundTripper struct {
	RoundTripFunc    func(req *http.Request) (*http.Response, error)
	RequestCount     int
	LastRequest      *http.Request
	RecordedRequests []*http.Request
}

// RoundTrip implements the http.RoundTripper interface
func (m *MockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	m.RequestCount++
	m.LastRequest = req
	m.RecordedRequests = append(m.RecordedRequests, req.Clone(req.Context()))

	if m.RoundTripFunc != nil {
		return m.RoundTripFunc(req)
	}

	// Default success response
	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(strings.NewReader(`{"success": true}`)),
		Header:     http.Header{},
	}, nil
}

// --- Error Handling Utilities ---

// createTestErrDecoder creates an ErrDecoder that tracks its calls
func createTestErrDecoder() (ErrDecoder, *bool) {
	called := false
	decoder := func(ctx context.Context, resp *resty.Response) error {
		called = true
		if resp.StatusCode() >= 400 {
			return NewGenericClientError("test", resp.StatusCode(), resp.Body())
		}
		return nil
	}
	return decoder, &called
}

// --- Driver Wrapper Utilities ---

// createDriverWrapper creates a DriverWrapper that tracks its calls
func createDriverWrapper(header, value string) (DriverWrapper, *bool) {
	called := false
	wrapper := func(client *resty.Client) *resty.Client {
		called = true
		return client.SetHeader(header, value)
	}
	return wrapper, &called
}

// --- Circuit Breaker Utilities ---

// resetHystrix clears all Hystrix command configurations and metrics
// Useful for test isolation between different test cases
func resetHystrix() {
	hystrix.Flush()
}

// withCircuitCommand adds circuit command to context for tests
func withCircuitCommand(ctx context.Context, command string) context.Context {
	return context.WithValue(ctx, circuitCommandKey, command)
}

// triggerCircuitOpen forces a circuit to open for testing circuit breaker recovery
func triggerCircuitOpen(command string) {
	// Configure the command with a low error threshold
	hystrix.ConfigureCommand(command, hystrix.CommandConfig{
		Timeout:                10,
		MaxConcurrentRequests:  1,
		ErrorPercentThreshold:  1, // Will open with just one error
		RequestVolumeThreshold: 1, // Only need one request to trigger
		SleepWindow:            100,
	})

	// Run a function that will time out to force the circuit to open
	_ = hystrix.Do(command, func() error {
		time.Sleep(20 * time.Millisecond) // Longer than the timeout
		return nil
	}, nil)
}

// --- Fallback Utilities ---

// createTestFallbackFunc creates a fallback function for testing with tracked calls
func createTestFallbackFunc(responseData interface{}) (fallbackFunc, *bool, *error) {
	called := false
	var passedError error

	fn := func(ctx context.Context, err error) (interface{}, error) {
		called = true
		passedError = err
		return responseData, nil
	}

	return fn, &called, &passedError
}

// createFailingFallbackFunc creates a fallback function that returns an error
func createFailingFallbackFunc() (fallbackFunc, *bool) {
	called := false

	fn := func(ctx context.Context, err error) (interface{}, error) {
		called = true
		return nil, err // Simply returns the passed error
	}

	return fn, &called
}
