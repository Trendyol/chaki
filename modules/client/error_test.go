package client

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestErrorResponse verifies that error responses are correctly created and handled
func TestErrorResponse(t *testing.T) {
	// Setup
	req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "https://example.com", nil)

	// Test cases for different error responses
	testCases := []struct {
		name           string
		statusCode     int
		body           string
		expectedStatus int
	}{
		{
			name:           "bad request error",
			statusCode:     http.StatusBadRequest,
			body:           `{"error": "Bad Request"}`,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "unauthorized error",
			statusCode:     http.StatusUnauthorized,
			body:           `{"error": "Unauthorized"}`,
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "internal server error",
			statusCode:     http.StatusInternalServerError,
			body:           `{"error": "Internal Server Error"}`,
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:           "default error message",
			statusCode:     http.StatusServiceUnavailable,
			body:           "", // Empty body to test default error message
			expectedStatus: http.StatusServiceUnavailable,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create error response using the utility function
			resp := createErrorResponse(req, tc.statusCode, tc.body)

			// Verify the response
			assert.NotNil(t, resp)
			assert.Equal(t, tc.expectedStatus, resp.StatusCode)
			assert.Equal(t, "application/json", resp.Header.Get("Content-Type"))
			assert.Equal(t, req, resp.Request)
		})
	}
}
