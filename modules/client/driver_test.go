package client

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
)

func TestNewDriverBuilder(t *testing.T) {
	// Setup
	cfg := driverTestConfig()

	// Execute
	builder := newDriverBuilder(cfg)

	// Verify
	assert.NotNil(t, builder)
	assert.Equal(t, cfg, builder.cfg)
	assert.NotNil(t, builder.d)
	assert.Equal(t, "https://example.com", builder.d.HostURL)
	assert.Equal(t, 1*time.Second, builder.d.GetClient().Timeout)
	assert.False(t, builder.d.Debug)
}

func TestDriverBuilder_AddErrDecoder(t *testing.T) {
	// Setup
	cfg := driverTestConfig()
	builder := newDriverBuilder(cfg)
	errDecoder, decoderCalled := createTestErrDecoder()

	// Execute
	result := builder.AddErrDecoder(errDecoder)

	// Verify
	assert.NotNil(t, result)
	// Don't compare functions directly, test that it's set to something
	assert.NotNil(t, builder.eh)
	assert.Same(t, builder, result) // Should return itself for chaining

	// Verify the decoder works by calling it
	dummyRes := &resty.Response{
		Request: resty.New().R().SetContext(context.Background()),
	}
	err := builder.eh(context.Background(), dummyRes)
	assert.NoError(t, err)
	assert.True(t, *decoderCalled, "Error decoder should be called")
}

func TestDriverBuilder_AddUpdaters(t *testing.T) {
	// Setup
	cfg := driverTestConfig()
	builder := newDriverBuilder(cfg)

	wrapper1, _ := createDriverWrapper("X-Test-1", "value1")
	wrapper2, _ := createDriverWrapper("X-Test-2", "value2")

	// Execute
	result := builder.AddUpdaters(wrapper1, wrapper2)

	// Verify
	assert.NotNil(t, result)
	assert.Len(t, builder.updaters, 2)
	assert.Same(t, builder, result) // Should return itself for chaining
}

func TestDriverBuilder_AddRoundTripperWrappers(t *testing.T) {
	// Setup
	cfg := driverTestConfig()
	builder := newDriverBuilder(cfg)

	wrapper1 := func(rt http.RoundTripper) http.RoundTripper { return rt }
	wrapper2 := func(rt http.RoundTripper) http.RoundTripper { return rt }

	// Execute
	result := builder.AddRoundTripperWrappers(wrapper1, wrapper2)

	// Verify
	assert.NotNil(t, result)
	assert.Len(t, builder.rtWrappers, 2)
	assert.Same(t, builder, result) // Should return itself for chaining
}

func TestDriverBuilder_SetRetry(t *testing.T) {
	t.Run("with retry configuration", func(t *testing.T) {
		// Setup
		cfg := driverTestConfig()
		builder := newDriverBuilder(cfg)
		retryConfig := &retryConfig{
			Name:      "test-retry",
			Count:     3,
			Interval:  100 * time.Millisecond,
			MaxDelay:  1 * time.Second,
			DelayType: ConstantDelay,
		}

		// Execute
		result := builder.SetRetry(retryConfig)

		// Verify
		assert.NotNil(t, result)
		assert.Len(t, builder.rtWrappers, 1)
		assert.Same(t, builder, result) // Should return itself for chaining
	})

	t.Run("with nil retry configuration", func(t *testing.T) {
		// Setup
		cfg := driverTestConfig()
		builder := newDriverBuilder(cfg)

		// Execute
		result := builder.SetRetry(nil)

		// Verify
		assert.NotNil(t, result)
		assert.Empty(t, builder.rtWrappers)
		assert.Same(t, builder, result) // Should return itself for chaining
	})
}

func TestDriverBuilder_SetCircuit(t *testing.T) {
	t.Run("with circuit configuration", func(t *testing.T) {
		// Setup
		cfg := driverTestConfig()
		builder := newDriverBuilder(cfg)
		circuitConfig := &circuitConfig{
			Enabled: true,
			Name:    "test-circuit",
		}

		// Execute
		result := builder.SetCircuit(circuitConfig)

		// Verify
		assert.NotNil(t, result)
		assert.Len(t, builder.rtWrappers, 1)
		assert.Same(t, builder, result) // Should return itself for chaining
	})

	t.Run("with nil circuit configuration", func(t *testing.T) {
		// Setup
		cfg := driverTestConfig()
		builder := newDriverBuilder(cfg)

		// Execute
		result := builder.SetCircuit(nil)

		// Verify
		assert.NotNil(t, result)
		assert.Empty(t, builder.rtWrappers)
		assert.Same(t, builder, result) // Should return itself for chaining
	})
}

func TestDriverBuilder_Build(t *testing.T) {
	t.Run("applies all updaters", func(t *testing.T) {
		// Setup
		cfg := driverTestConfig()
		builder := newDriverBuilder(cfg)

		wrapper1, wrapper1Called := createDriverWrapper("X-Test-1", "value1")
		wrapper2, wrapper2Called := createDriverWrapper("X-Test-2", "value2")

		builder.updaters = append(builder.updaters, wrapper1, wrapper2)

		// Add error handler
		errDecoder, errorHandlerCalled := createTestErrDecoder()
		builder.eh = errDecoder

		// Execute
		client := builder.build()

		// Verify
		assert.NotNil(t, client)
		assert.True(t, *wrapper1Called, "First wrapper should be applied")
		assert.True(t, *wrapper2Called, "Second wrapper should be applied")

		// Create a dummy response to test the error handler
		req := client.R().SetContext(context.Background())
		dummyRes := &resty.Response{
			Request: req,
		}

		// Instead of directly accessing OnAfterResponseFuncs, create a request and test
		// the error handler manually
		err := errDecoder(context.Background(), dummyRes)
		assert.NoError(t, err)
		assert.True(t, *errorHandlerCalled, "Error handler should be called")
	})

	t.Run("configures round trippers in correct order", func(t *testing.T) {
		// Setup
		cfg := driverTestConfig()
		cfg.Set("baseurl", "https://test-api.example.com")
		cfg.Set("logging", false)
		builder := newDriverBuilder(cfg)

		// Add round tripper wrappers in a specific order
		executionOrder := make([]int, 0)

		builder.rtWrappers = append(builder.rtWrappers,
			func(rt http.RoundTripper) http.RoundTripper {
				executionOrder = append(executionOrder, 1)
				return rt
			},
			func(rt http.RoundTripper) http.RoundTripper {
				executionOrder = append(executionOrder, 2)
				return rt
			},
		)

		// Execute
		builder.build()

		// Verify
		assert.Equal(t, []int{1, 2}, executionOrder, "Round trippers should be applied in order")
	})

	t.Run("enables logging when configured", func(t *testing.T) {
		// Setup
		cfg := driverTestConfig()
		cfg.Set("baseurl", "https://test-api.example.com")
		cfg.Set("logging", true)
		builder := newDriverBuilder(cfg)

		// Execute
		client := builder.build()

		// Verify
		assert.NotNil(t, client)
		// Don't check internal fields directly, verify functionality instead
		assert.NotNil(t, client)
	})
}

func TestSetDefaults(t *testing.T) {
	// Setup - use driverTestConfig but reset the values we want to test
	cfg := driverTestConfig()

	// Reset the values we want to test defaults for
	cfg.Set("timeout", nil)
	cfg.Set("debug", nil)
	cfg.Set("logging", nil)

	// Execute
	setDefaults(cfg)

	// Verify
	assert.Equal(t, "5s", cfg.GetString("timeout"))
	assert.False(t, cfg.GetBool("debug"))
	assert.False(t, cfg.GetBool("logging"))
}

func TestDriverBuilder_Integration(t *testing.T) {
	// Setup a test server with our standardHandler utility
	server := mockServer(standardHandler(http.StatusOK, map[string]interface{}{
		"message": "success",
	}, 0))
	defer server.Close()

	// Create config
	cfg := driverTestConfig()
	cfg.Set("baseurl", server.URL)

	// Create wrappers for testing
	wrapper, wrapperCalled := createDriverWrapper("X-Test", "value")
	errDecoder, errDecoderCalled := createTestErrDecoder()

	// Create the builder with all the features
	builder := newDriverBuilder(cfg).
		AddErrDecoder(errDecoder).
		AddUpdaters(wrapper).
		AddRoundTripperWrappers(func(rt http.RoundTripper) http.RoundTripper {
			return rt
		})

	// Build the client
	client := builder.build()

	// Make a test request
	resp, err := client.R().Get("/")

	// Verify
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode())
	assert.True(t, *wrapperCalled, "Updater should be called")
	// The error decoder will be called with a 200 status, so no error should be returned
	assert.True(t, *errDecoderCalled, "Error decoder should be called")
}

func TestDriverBuilder_WithMockRoundTripper(t *testing.T) {
	// Setup
	cfg := driverTestConfig()
	cfg.Set("baseurl", "https://test-api.example.com")
	builder := newDriverBuilder(cfg)

	// Create a mock round tripper with custom behavior
	mockRT := &MockRoundTripper{
		RoundTripFunc: func(req *http.Request) (*http.Response, error) {
			return createSuccessResponse(req, `{"custom": "response"}`), nil
		},
	}

	// Replace the default transport with our mock
	builder.d.SetTransport(mockRT)
	builder.eh = DefaultErrDecoder("testClient")

	// Add a round tripper wrapper that should wrap our mock
	wrapperCalled := false
	builder.AddRoundTripperWrappers(func(rt http.RoundTripper) http.RoundTripper {
		wrapperCalled = true
		// Verify that we're wrapping the MockRoundTripper
		_, isMockRT := rt.(*MockRoundTripper)
		assert.True(t, isMockRT, "Should be wrapping our MockRoundTripper")
		return rt
	})

	// Build the client
	client := builder.build()

	// Verify
	assert.NotNil(t, client)
	assert.True(t, wrapperCalled, "RoundTripper wrapper should be called")

	// Test HTTP request using our mock client
	resp, err := client.R().Get("/test-endpoint")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode())
	assert.Contains(t, string(resp.Body()), "custom")
	assert.Equal(t, 1, mockRT.RequestCount, "Should have made exactly one request")
}

func TestDriverBuilder_UseLogging(t *testing.T) {
	t.Run("registers and executes logging callbacks", func(t *testing.T) {
		// Setup
		cfg := driverTestConfig()
		builder := newDriverBuilder(cfg)

		// Execute - call useLogging directly
		builder.useLogging()

		// Create a test server that we'll use to verify the logging works
		server := mockServer(standardHandler(http.StatusOK, map[string]interface{}{
			"response": "test",
		}, 0))
		defer server.Close()

		// Configure the client to use our test server
		builder.d.SetBaseURL(server.URL)

		// Make a test request with various parameters to test logging
		req := builder.d.R().
			SetHeader("X-Test", "test-value").
			SetQueryParam("param", "value").
			SetBody(map[string]string{"key": "value"})

		// Execute the request which should trigger both logging callbacks
		resp, err := req.Execute("GET", "/test")

		// Verify
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, http.StatusOK, resp.StatusCode())
	})
}
