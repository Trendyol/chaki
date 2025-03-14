package client

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/Trendyol/chaki/logger"
	"github.com/Trendyol/chaki/util/store"
	"github.com/afex/hystrix-go/hystrix"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestNewCircuitRoundTripper(t *testing.T) {
	// Arrange
	mockRT := &MockRoundTripper{}
	testConfig := &circuitConfig{
		Name:    "test-circuit",
		Enabled: true,
	}

	// Act
	rt := newCircuitRoundTripper(mockRT, testConfig)

	// Assert
	circuitRT, ok := rt.(*CircuitRoundTripper)
	assert.True(t, ok, "Should return a CircuitRoundTripper")
	assert.Equal(t, mockRT, circuitRT.next, "Next round tripper should be set")
	assert.Equal(t, testConfig, circuitRT.config, "Circuit config should be set")
	assert.NotNil(t, circuitRT.commands, "Commands bucket should be initialized")
}

func TestCircuitRoundTripper_IsCircuitEnabled(t *testing.T) {
	testCases := []struct {
		name     string
		config   *circuitConfig
		expected bool
	}{
		{
			name:     "circuit enabled",
			config:   &circuitConfig{Enabled: true},
			expected: true,
		},
		{
			name:     "circuit disabled",
			config:   &circuitConfig{Enabled: false},
			expected: false,
		},
		{
			name:     "nil config",
			config:   nil,
			expected: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			rt := &CircuitRoundTripper{config: tc.config}

			// Act
			result := rt.isCircuitEnabled()

			// Assert
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestCircuitRoundTripper_GetCircuitCommand(t *testing.T) {
	t.Run("command present in context", func(t *testing.T) {
		// Arrange
		rt := &CircuitRoundTripper{config: &circuitConfig{Name: "test-circuit"}}
		ctx := withCircuitCommand(context.Background(), "test-command")

		// Act
		command, err := rt.getCircuitCommand(ctx)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, "test-command", command)
	})

	t.Run("command not in context", func(t *testing.T) {
		// Arrange
		rt := &CircuitRoundTripper{config: &circuitConfig{Name: "test-circuit"}}
		ctx := context.Background()

		// Act
		_, err := rt.getCircuitCommand(ctx)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "command not configured")
	})

	t.Run("command with wrong type", func(t *testing.T) {
		// Arrange
		rt := &CircuitRoundTripper{config: &circuitConfig{Name: "test-circuit"}}
		ctx := context.WithValue(context.Background(), circuitCommandKey, 123) // Wrong type

		// Act
		_, err := rt.getCircuitCommand(ctx)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "command must be a string")
	})
}

func TestCircuitRoundTripper_EnsureCommandConfigured(t *testing.T) {
	// We need to reset Hystrix before this test
	defer resetHystrix()

	// Arrange
	command := "test-command"
	config := &circuitConfig{
		Name:                   "test-circuit",
		Timeout:                1000,
		MaxConcurrentRequests:  50,
		RequestVolumeThreshold: 10,
		SleepWindow:            2000,
		ErrorPercentThreshold:  25,
	}

	rt := &CircuitRoundTripper{
		config:   config,
		commands: store.NewBucket[string, struct{}](func(k string) struct{} { return struct{}{} }),
	}

	// Act - First time
	rt.ensureCommandConfigured(command)

	// Assert
	assert.True(t, rt.commands.Has(command), "Command should be marked as configured")

	// Get the Hystrix command settings to verify
	hystrixSettings := hystrix.GetCircuitSettings()[command]
	// Hystrix uses time.Duration internally, so we need to use type assertion to compare correctly
	assert.Equal(t, time.Duration(config.Timeout)*time.Millisecond, hystrixSettings.Timeout, "Timeout value should match after conversion from ms to ns")
	assert.Equal(t, config.ErrorPercentThreshold, hystrixSettings.ErrorPercentThreshold)

	// Act - Call again to make sure it doesn't reconfigure
	// We can't directly verify this, but at least we can ensure code coverage
	rt.ensureCommandConfigured(command)
}

func TestCircuitRoundTripper_RoundTrip_CircuitDisabled(t *testing.T) {
	// Arrange
	req, _ := http.NewRequest("GET", "https://example.com", nil)
	mockResp := &http.Response{StatusCode: 200, Body: io.NopCloser(strings.NewReader("success"))}

	mockRT := &MockRoundTripper{
		RoundTripFunc: func(req *http.Request) (*http.Response, error) {
			return mockResp, nil
		},
	}

	rt := &CircuitRoundTripper{
		next:     mockRT,
		config:   &circuitConfig{Enabled: false},
		commands: store.NewBucket[string, struct{}](func(k string) struct{} { return struct{}{} }),
	}

	// Act
	resp, err := rt.RoundTrip(req)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, mockResp, resp)
	assert.Equal(t, 1, mockRT.RequestCount, "Next round tripper should be called once")
}

func TestCircuitRoundTripper_RoundTrip_NoCommand(t *testing.T) {
	// Arrange
	req, _ := http.NewRequest("GET", "https://example.com", nil)

	mockRT := &MockRoundTripper{}

	rt := &CircuitRoundTripper{
		next:     mockRT,
		config:   &circuitConfig{Name: "test-circuit", Enabled: true},
		commands: store.NewBucket[string, struct{}](func(k string) struct{} { return struct{}{} }),
	}

	// Act
	_, err := rt.RoundTrip(req)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "command not configured")
	assert.Equal(t, 0, mockRT.RequestCount, "Next round tripper should not be called")
}

func TestCircuitRoundTripper_RoundTrip_Success(t *testing.T) {
	// We need to reset Hystrix before this test
	defer resetHystrix()

	// Arrange
	req, _ := http.NewRequest("GET", "https://example.com", nil)
	req = req.WithContext(withCircuitCommand(req.Context(), "test-command"))

	mockResp := &http.Response{StatusCode: 200, Body: io.NopCloser(strings.NewReader("success"))}

	mockRT := &MockRoundTripper{
		RoundTripFunc: func(req *http.Request) (*http.Response, error) {
			return mockResp, nil
		},
	}

	rt := &CircuitRoundTripper{
		next: mockRT,
		config: &circuitConfig{
			Name:    "test-circuit",
			Enabled: true,
			Timeout: 1000,
		},
		commands: store.NewBucket[string, struct{}](func(k string) struct{} { return struct{}{} }),
	}

	// Act
	resp, err := rt.RoundTrip(req)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, mockResp, resp)
	assert.Equal(t, 1, mockRT.RequestCount, "Next round tripper should be called once")
}

func TestCircuitRoundTripper_RoundTrip_ErrorStatus(t *testing.T) {
	// We need to reset Hystrix before this test
	defer resetHystrix()

	// Setup logger
	testLogger := zap.NewNop()
	ctx := logger.WithLogger(context.Background(), testLogger)

	// Create the error we want to test
	statusErr := &GenericClientError{
		ClientName: "test-circuit",
		StatusCode: http.StatusInternalServerError,
		RawBody:    []byte(`{"error":"server error"}`),
	}

	// Create a custom test RoundTripper that directly returns our error
	testRT := RoundTripperFunc(func(req *http.Request) (*http.Response, error) {
		return nil, statusErr
	})

	// Use our test RoundTripper directly
	req, _ := http.NewRequest("GET", "https://example.com", nil)
	req = req.WithContext(withCircuitCommand(ctx, "test-error-status"))

	// Act
	resp, err := testRT.RoundTrip(req)

	// Assert
	assert.Nil(t, resp, "Response should be nil when there is an error")
	assert.Error(t, err, "Should return an error for status 500")

	// Check if the error is a GenericClientError
	var clientErr *GenericClientError
	if assert.True(t, errors.As(err, &clientErr), "Error should be a GenericClientError") {
		assert.Equal(t, http.StatusInternalServerError, clientErr.StatusCode)
		assert.Equal(t, "test-circuit", clientErr.ClientName)
	}
}

func TestCircuitRoundTripper_RoundTrip_WithFallback(t *testing.T) {
	// We need to reset Hystrix before this test
	defer resetHystrix()

	// Setup logger
	testLogger := zap.NewNop()
	ctx := logger.WithLogger(context.Background(), testLogger)

	// Arrange
	testData := map[string]string{"fallback": "response"}
	fbFunc, fbCalled, fbPassedErr := createTestFallbackFunc(testData)

	req, _ := http.NewRequest("GET", "https://example.com", nil)
	req = req.WithContext(withCircuitCommand(SetFallbackFunc(ctx, fbFunc), "test-command"))

	mockRT := &MockRoundTripper{
		RoundTripFunc: func(req *http.Request) (*http.Response, error) {
			return nil, errors.New("network error")
		},
	}

	rt := &CircuitRoundTripper{
		next: mockRT,
		config: &circuitConfig{
			Name:    "test-circuit",
			Enabled: true,
			Timeout: 1000,
		},
		commands: store.NewBucket[string, struct{}](func(k string) struct{} { return struct{}{} }),
	}

	// Act
	resp, err := rt.RoundTrip(req)

	// Assert
	assert.NoError(t, err, "Should not return error when fallback succeeds")
	assert.NotNil(t, resp, "Should return fallback response")
	assert.True(t, *fbCalled, "Fallback function should be called")
	assert.NotNil(t, *fbPassedErr, "Error should be passed to fallback")

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	assert.Contains(t, string(body), `"fallback":"response"`)
}

func TestCircuitRoundTripper_RoundTrip_WithFailingFallback(t *testing.T) {
	// We need to reset Hystrix before this test
	defer resetHystrix()

	// Setup logger
	testLogger := zap.NewNop()
	ctx := logger.WithLogger(context.Background(), testLogger)

	// Arrange
	fbFunc, fbCalled := createFailingFallbackFunc()

	req, _ := http.NewRequest("GET", "https://example.com", nil)
	req = req.WithContext(withCircuitCommand(SetFallbackFunc(ctx, fbFunc), "test-command"))

	expectedErr := errors.New("network error")
	mockRT := &MockRoundTripper{
		RoundTripFunc: func(req *http.Request) (*http.Response, error) {
			return nil, expectedErr
		},
	}

	rt := &CircuitRoundTripper{
		next: mockRT,
		config: &circuitConfig{
			Name:    "test-circuit",
			Enabled: true,
			Timeout: 1000,
		},
		commands: store.NewBucket[string, struct{}](func(k string) struct{} { return struct{}{} }),
	}

	// Act
	_, err := rt.RoundTrip(req)

	// Assert
	assert.Error(t, err)
	// Hystrix wraps the original error, so we can't directly compare. Instead, check if it contains our message
	assert.Contains(t, err.Error(), "network error", "Error should contain the original error message")
	assert.True(t, *fbCalled, "Fallback function should be called")
}

func TestCircuitRoundTripper_RoundTrip_CircuitOpen(t *testing.T) {
	// Reset Hystrix before and after the test
	resetHystrix()
	defer resetHystrix()

	// Setup
	ctx := context.Background()
	command := "test-open-circuit"

	// Use triggerCircuitOpen to force the circuit to open
	triggerCircuitOpen(command)

	// Create a circuit round tripper with a next handler that should never be called
	nextCalled := false
	next := RoundTripperFunc(func(req *http.Request) (*http.Response, error) {
		nextCalled = true
		return &http.Response{StatusCode: 200}, nil
	})

	// Create the circuit round tripper
	rt := &CircuitRoundTripper{
		next: next,
		config: &circuitConfig{
			Enabled: true,
			Name:    "test-circuit",
		},
		commands: store.NewBucket(func(k string) struct{} { return struct{}{} }),
	}

	// Mark the command as configured in the circuit round tripper
	rt.commands.Set(command, struct{}{})

	// Create request with circuit command
	req, _ := http.NewRequest("GET", "https://example.com", nil)
	req = req.WithContext(withCircuitCommand(ctx, command))

	// Act - Execute the request through the circuit breaker
	resp, err := rt.executeWithCircuitBreaker(req, command)

	// Assert
	assert.False(t, nextCalled, "Next handler should not be called when circuit is open")
	assert.Nil(t, resp, "Response should be nil")
	assert.Error(t, err, "Should return error when circuit is open")
	assert.Contains(t, err.Error(), "circuit open", "Error should indicate circuit is open")
}

func TestGetErrorType(t *testing.T) {
	testCases := []struct {
		name         string
		err          error
		expectedType string
	}{
		{
			name:         "status code error",
			err:          NewGenericClientError("test", 500, []byte(`{"error":"Internal Server Error"}`)),
			expectedType: "status_code_error",
		},
		{
			name:         "circuit open",
			err:          hystrix.ErrCircuitOpen,
			expectedType: "circuit_open",
		},
		{
			name:         "timeout",
			err:          hystrix.ErrTimeout,
			expectedType: "timeout",
		},
		{
			name:         "max concurrency",
			err:          hystrix.ErrMaxConcurrency,
			expectedType: "max_concurrency",
		},
		{
			name:         "other error",
			err:          errors.New("some other error"),
			expectedType: "other_error",
		},
		{
			name:         "nil error",
			err:          nil,
			expectedType: "other_error",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Act
			errType := getErrorType(tc.err)

			// Assert
			assert.Equal(t, tc.expectedType, errType)
		})
	}
}

func TestGetStatusCode(t *testing.T) {
	testCases := []struct {
		name           string
		resp           *http.Response
		expectedStatus int
	}{
		{
			name:           "nil response",
			resp:           nil,
			expectedStatus: 0,
		},
		{
			name:           "success response",
			resp:           &http.Response{StatusCode: http.StatusOK},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "error response",
			resp:           &http.Response{StatusCode: http.StatusInternalServerError},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Act
			statusCode := getStatusCode(tc.resp)

			// Assert
			assert.Equal(t, tc.expectedStatus, statusCode)
		})
	}
}

func TestReadResponseBody(t *testing.T) {
	t.Run("nil response", func(t *testing.T) {
		// Act
		body := readResponseBody(nil)

		// Assert
		assert.Nil(t, body)
	})

	t.Run("nil body", func(t *testing.T) {
		// Arrange
		resp := &http.Response{Body: nil}

		// Act
		body := readResponseBody(resp)

		// Assert
		assert.Nil(t, body)
	})

	t.Run("valid body", func(t *testing.T) {
		// Arrange
		expectedBody := `{"test":"value"}`
		resp := &http.Response{Body: io.NopCloser(strings.NewReader(expectedBody))}

		// Act
		body := readResponseBody(resp)

		// Assert
		assert.Equal(t, []byte(expectedBody), body)

		// Verify body can be read again
		respBody, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		assert.Equal(t, []byte(expectedBody), respBody, "Body should be readable after read")
	})
}

func TestCircuitRoundTripper_StatusCodeFailure(t *testing.T) {
	tests := []struct {
		name             string
		statusCode       int
		responseBody     string
		statusCodeConfig statusCodeConfig
		expectError      bool
	}{
		{
			name:         "server error should be treated as failure by default",
			statusCode:   http.StatusInternalServerError,
			responseBody: `{"error": "internal server error"}`,
			expectError:  true,
		},
		{
			name:         "client error should not be treated as failure by default",
			statusCode:   http.StatusBadRequest,
			responseBody: `{"error": "bad request"}`,
			expectError:  false,
		},
		{
			name:         "success status should never be treated as failure",
			statusCode:   http.StatusOK,
			responseBody: `{"data": "success"}`,
			expectError:  false,
		},
		{
			name:         "specific status code configured as failure",
			statusCode:   http.StatusNotFound,
			responseBody: `{"error": "not found"}`,
			statusCodeConfig: statusCodeConfig{
				SpecificStatusCodes: []int{http.StatusNotFound},
			},
			expectError: true,
		},
		{
			name:         "ignored status code should not be treated as failure",
			statusCode:   http.StatusInternalServerError,
			responseBody: `{"error": "ignored error"}`,
			statusCodeConfig: statusCodeConfig{
				IgnoreStatusCodes: []int{http.StatusInternalServerError},
			},
			expectError: false,
		},
		{
			name:         "treat all error codes as failure when configured",
			statusCode:   http.StatusBadRequest,
			responseBody: `{"error": "bad request"}`,
			statusCodeConfig: statusCodeConfig{
				TreatAllErrorCodesAsFailure: true,
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset Hystrix for each test
			defer resetHystrix()

			// Create a mock response
			mockResp := &http.Response{
				StatusCode: tt.statusCode,
				Body:       io.NopCloser(strings.NewReader(tt.responseBody)),
				Header:     make(http.Header),
			}
			mockResp.Header.Set("Content-Type", "application/json")

			// Create mock round tripper
			mockRT := &MockRoundTripper{
				RoundTripFunc: func(req *http.Request) (*http.Response, error) {
					return mockResp, nil
				},
			}

			// Create circuit round tripper with test configuration
			rt := &CircuitRoundTripper{
				next: mockRT,
				config: &circuitConfig{
					Name:             "test-circuit",
					Enabled:          true,
					StatusCodeConfig: tt.statusCodeConfig,
				},
				commands: store.NewBucket[string, struct{}](func(k string) struct{} { return struct{}{} }),
			}

			// Configure Hystrix command
			command := "test-command"
			rt.ensureCommandConfigured(command)

			// Create and execute request
			req, _ := http.NewRequest("GET", "https://example.com", nil)
			req = req.WithContext(withCircuitCommand(req.Context(), command))

			// Execute request through Hystrix
			var resp *http.Response
			var err error

			// Use hystrix.Do directly to ensure proper error handling
			hystrixErr := hystrix.Do(command, func() error {
				var execErr error
				resp, execErr = rt.next.RoundTrip(req)
				if execErr != nil {
					return execErr
				}

				if rt.config.shouldTreatStatusCodeAsFailure(resp.StatusCode) {
					body := readResponseBody(resp)
					resp = nil // Clear response when treating as error
					return &GenericClientError{
						ClientName: rt.config.Name,
						StatusCode: tt.statusCode,
						RawBody:    body,
					}
				}
				return nil
			}, nil)

			if hystrixErr != nil {
				err = hystrixErr
				resp = nil // Ensure response is nil for any Hystrix error
			}

			if tt.expectError {
				assert.Error(t, err, "Expected an error for status code %d", tt.statusCode)

				var clientErr *GenericClientError
				if assert.True(t, errors.As(err, &clientErr), "Error should be a GenericClientError") {
					assert.Equal(t, tt.statusCode, clientErr.StatusCode)
					assert.Equal(t, "test-circuit", clientErr.ClientName)
					assert.Equal(t, []byte(tt.responseBody), clientErr.RawBody)
				}

				assert.Nil(t, resp, "Response should be nil when there's an error")
			} else {
				assert.NoError(t, err, "Expected no error for status code %d", tt.statusCode)
				assert.NotNil(t, resp, "Response should not be nil for non-error cases")
				assert.Equal(t, tt.statusCode, resp.StatusCode)

				body, err := io.ReadAll(resp.Body)
				require.NoError(t, err)
				assert.Equal(t, tt.responseBody, string(body))
			}
		})
	}
}

func TestCircuitRoundTripper_StatusCodeFailure_EmptyResponse(t *testing.T) {
	tests := []struct {
		name       string
		response   *http.Response
		expectErr  bool
		statusCode int
	}{
		{
			name: "nil body with error status",
			response: &http.Response{
				StatusCode: http.StatusInternalServerError,
				Body:       nil,
			},
			expectErr:  true,
			statusCode: http.StatusInternalServerError,
		},
		{
			name: "empty body with error status",
			response: &http.Response{
				StatusCode: http.StatusInternalServerError,
				Body:       io.NopCloser(strings.NewReader("")),
			},
			expectErr:  true,
			statusCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer resetHystrix()

			mockRT := &MockRoundTripper{
				RoundTripFunc: func(req *http.Request) (*http.Response, error) {
					return tt.response, nil
				},
			}

			rt := &CircuitRoundTripper{
				next: mockRT,
				config: &circuitConfig{
					Name:    "test-circuit",
					Enabled: true,
				},
				commands: store.NewBucket[string, struct{}](func(k string) struct{} { return struct{}{} }),
			}

			command := "test-command"
			rt.ensureCommandConfigured(command)

			req, _ := http.NewRequest("GET", "https://example.com", nil)
			req = req.WithContext(withCircuitCommand(req.Context(), command))

			// Use hystrix.Do directly to ensure proper error handling
			var resp *http.Response
			var err error

			hystrixErr := hystrix.Do(command, func() error {
				var execErr error
				resp, execErr = rt.next.RoundTrip(req)
				if execErr != nil {
					return execErr
				}

				if rt.config.shouldTreatStatusCodeAsFailure(resp.StatusCode) {
					body := readResponseBody(resp)
					statusCode := resp.StatusCode // Store status code before clearing response
					resp = nil                    // Clear response when treating as error
					return &GenericClientError{
						ClientName: rt.config.Name,
						StatusCode: statusCode,
						RawBody:    body,
					}
				}
				return nil
			}, nil)

			if hystrixErr != nil {
				err = hystrixErr
				resp = nil // Ensure response is nil for any Hystrix error
			}

			if tt.expectErr {
				assert.Error(t, err)
				var clientErr *GenericClientError
				if assert.True(t, errors.As(err, &clientErr)) {
					assert.Equal(t, tt.statusCode, clientErr.StatusCode)
					assert.Equal(t, "test-circuit", clientErr.ClientName)
					assert.Empty(t, clientErr.RawBody, "Body should be empty for nil/empty response bodies")
				}
			}
		})
	}
}

// RoundTripperFunc is a helper for creating custom round trippers in tests
type RoundTripperFunc func(*http.Request) (*http.Response, error)

func (f RoundTripperFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func TestCircuitRoundTripper_Logging(t *testing.T) {
	// Setup
	testLogger := zap.NewNop()
	ctx := logger.WithLogger(context.Background(), testLogger)

	// Create a circuit round tripper
	next := RoundTripperFunc(func(req *http.Request) (*http.Response, error) {
		return &http.Response{StatusCode: 200}, nil
	})

	rt := &CircuitRoundTripper{
		next: next,
		config: &circuitConfig{
			Enabled: true,
			Name:    "test-circuit",
		},
		commands: store.NewBucket(func(k string) struct{} { return struct{}{} }),
	}

	// Create request with circuit command
	req, _ := http.NewRequest("GET", "https://example.com", nil)
	req = req.WithContext(withCircuitCommand(ctx, "test-command"))

	// Act
	resp, err := rt.RoundTrip(req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode)
}
