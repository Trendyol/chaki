package client

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type MockRetryRoundTripper struct {
	callCount    int
	responses    []*http.Response
	errors       []error
	requestDelay time.Duration
}

func (m *MockRetryRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	if m.requestDelay > 0 {
		time.Sleep(m.requestDelay)
	}

	if m.callCount < len(m.responses) {
		resp := m.responses[m.callCount]
		err := m.errors[m.callCount]
		m.callCount++
		return resp, err
	}

	return nil, errors.New("unexpected call to RoundTrip")
}

func TestNewRetryRoundTripper(t *testing.T) {
	// Arrange
	mockRT := &MockRetryRoundTripper{}
	cfg := &retryConfig{
		Name:      "test",
		Count:     3,
		Interval:  100 * time.Millisecond,
		MaxDelay:  1 * time.Second,
		DelayType: ConstantDelay,
	}

	// Act
	rt := newRetryRoundTripper(mockRT, cfg)

	// Assert
	assert.NotNil(t, rt)
	retryRT, ok := rt.(*RetryRoundTripper)
	assert.True(t, ok)
	assert.Equal(t, mockRT, retryRT.next)
	assert.Equal(t, cfg, retryRT.cfg)
}

func TestRetryRoundTripper_RoundTrip(t *testing.T) {
	testCases := []struct {
		name          string
		config        *retryConfig
		responses     []*http.Response
		errors        []error
		expectedCalls int
		expectedError bool
		timeout       time.Duration
	}{
		{
			name: "success on first try",
			config: &retryConfig{
				Count:     3,
				Interval:  10 * time.Millisecond,
				DelayType: ConstantDelay,
			},
			responses:     []*http.Response{{StatusCode: 200}},
			errors:        []error{nil},
			expectedCalls: 1,
			expectedError: false,
		},
		{
			name: "success after retries with constant delay",
			config: &retryConfig{
				Count:     3,
				Interval:  10 * time.Millisecond,
				DelayType: ConstantDelay,
			},
			responses:     []*http.Response{nil, nil, {StatusCode: 200}},
			errors:        []error{errors.New("error1"), errors.New("error2"), nil},
			expectedCalls: 3,
			expectedError: false,
		},
		{
			name: "success after retries with exponential delay",
			config: &retryConfig{
				Count:     3,
				Interval:  10 * time.Millisecond,
				MaxDelay:  100 * time.Millisecond,
				DelayType: ExponentialDelay,
			},
			responses:     []*http.Response{nil, nil, {StatusCode: 200}},
			errors:        []error{errors.New("error1"), errors.New("error2"), nil},
			expectedCalls: 3,
			expectedError: false,
		},
		{
			name: "failure after all retries",
			config: &retryConfig{
				Count:     2,
				Interval:  10 * time.Millisecond,
				DelayType: ConstantDelay,
			},
			responses:     []*http.Response{nil, nil, nil},
			errors:        []error{errors.New("error1"), errors.New("error2"), errors.New("error3")},
			expectedCalls: 3,
			expectedError: true,
		},
		{
			name:          "no retry when config is nil",
			config:        nil,
			responses:     []*http.Response{nil},
			errors:        []error{errors.New("error")},
			expectedCalls: 1,
			expectedError: true,
		},
		{
			name: "no retry when count is 0",
			config: &retryConfig{
				Count:     0,
				Interval:  10 * time.Millisecond,
				DelayType: ConstantDelay,
			},
			responses:     []*http.Response{nil},
			errors:        []error{errors.New("error")},
			expectedCalls: 1,
			expectedError: true,
		},
		{
			name: "respect context cancellation",
			config: &retryConfig{
				Count:     3,
				Interval:  100 * time.Millisecond,
				DelayType: ConstantDelay,
			},
			responses:     []*http.Response{nil},
			errors:        []error{errors.New("error")},
			expectedCalls: 1,
			expectedError: true,
			timeout:       50 * time.Millisecond,
		},
		{
			name: "respect max delay with exponential backoff",
			config: &retryConfig{
				Count:     3,
				Interval:  10 * time.Millisecond,
				MaxDelay:  15 * time.Millisecond,
				DelayType: ExponentialDelay,
			},
			responses:     []*http.Response{nil, nil, {StatusCode: 200}},
			errors:        []error{errors.New("error1"), errors.New("error2"), nil},
			expectedCalls: 3,
			expectedError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			mockRT := &MockRetryRoundTripper{
				responses: tc.responses,
				errors:    tc.errors,
			}

			rt := &RetryRoundTripper{
				next: mockRT,
				cfg:  tc.config,
			}

			ctx := context.Background()
			if tc.timeout > 0 {
				var cancel context.CancelFunc
				ctx, cancel = context.WithTimeout(ctx, tc.timeout)
				defer cancel()
			}

			req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://example.com", nil)
			require.NoError(t, err)

			// Act
			resp, err := rt.RoundTrip(req)

			// Assert
			assert.Equal(t, tc.expectedCalls, mockRT.callCount)
			if tc.expectedError {
				assert.Error(t, err)
				if tc.timeout > 0 {
					assert.ErrorIs(t, err, context.DeadlineExceeded)
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				assert.Equal(t, 200, resp.StatusCode)
			}
		})
	}
}

func TestRetryRoundTripper_ExponentialBackoff(t *testing.T) {
	// Arrange
	cfg := &retryConfig{
		Count:     3,
		Interval:  10 * time.Millisecond,
		MaxDelay:  50 * time.Millisecond,
		DelayType: ExponentialDelay,
	}

	mockRT := &MockRetryRoundTripper{
		responses:    []*http.Response{nil, nil, {StatusCode: 200}},
		errors:       []error{errors.New("error1"), errors.New("error2"), nil},
		requestDelay: 1 * time.Millisecond, // Small delay to make timing more realistic
	}

	rt := &RetryRoundTripper{
		next: mockRT,
		cfg:  cfg,
	}

	req, err := http.NewRequest(http.MethodGet, "http://example.com", nil)
	require.NoError(t, err)

	// Act
	start := time.Now()
	resp, err := rt.RoundTrip(req)
	duration := time.Since(start)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, 3, mockRT.callCount)

	// The actual duration should be greater than the sum of the minimum delays
	// First retry: 10ms
	// Second retry: 20ms + jitter
	minExpectedDuration := 30 * time.Millisecond
	assert.Greater(t, duration, minExpectedDuration)

	// But should not exceed maximum delay plus some buffer for execution time
	maxExpectedDuration := cfg.MaxDelay + (100 * time.Millisecond)
	assert.Less(t, duration, maxExpectedDuration)
}
