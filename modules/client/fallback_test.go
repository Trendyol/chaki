package client

import (
	"context"
	"errors"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSetFallbackFunc(t *testing.T) {
	// Arrange
	ctx := context.Background()
	var called bool
	testFn := func(ctx context.Context, err error) (interface{}, error) {
		called = true
		return nil, nil
	}

	// Act
	ctxWithFallback := SetFallbackFunc(ctx, testFn)

	// Assert
	assert.NotEqual(t, ctx, ctxWithFallback, "Context should be different after setting fallback")

	// Verify we can retrieve the function from context
	fnFromCtx, ok := ctxWithFallback.Value(circuitFallbackKey).(fallbackFunc)
	assert.True(t, ok, "Should be able to retrieve fallback function from context")

	// Call the retrieved function to verify it's our function
	_, _ = fnFromCtx(ctx, nil)
	assert.True(t, called, "Retrieved function should call our test function")
}

func TestNewOrDefaultFallbackHandler(t *testing.T) {
	t.Run("with custom fallback function", func(t *testing.T) {
		// Arrange
		testData := map[string]string{"key": "value"}
		fn, called, _ := createTestFallbackFunc(testData)
		ctx := SetFallbackFunc(context.Background(), fn)

		// Act
		handler := newOrDefaultFallbackHandler(ctx)

		// Assert - Test function behavior rather than comparing functions directly
		// Call both functions with same input and compare results
		testErr := errors.New("test error")
		result1, err1 := handler.fn(context.Background(), testErr)
		result2, err2 := fn(context.Background(), testErr)

		assert.Equal(t, result1, result2, "Handler function should return same result as original")
		assert.Equal(t, err1, err2, "Handler function should return same error as original")
		assert.True(t, *called, "Function should be called")
	})

	t.Run("with default fallback function", func(t *testing.T) {
		// Arrange
		ctx := context.Background()

		// Act
		handler := newOrDefaultFallbackHandler(ctx)

		// Execute to verify it's the default function
		testError := errors.New("test error")
		resp, err := handler.fn(ctx, testError)

		// Assert
		assert.Nil(t, resp, "Default function should return nil response")
		assert.Equal(t, testError, err, "Default function should return the original error")
	})
}

func TestDefaultFallbackFn(t *testing.T) {
	// Arrange
	testError := errors.New("test error")
	ctx := context.Background()

	// Act
	resp, err := defaultFallbackFn(ctx, testError)

	// Assert
	assert.Nil(t, resp, "Default fallback function should return nil response")
	assert.Equal(t, testError, err, "Default fallback function should return original error")
}

func TestFallbackHandler_Handle(t *testing.T) {
	t.Run("successful fallback", func(t *testing.T) {
		// Arrange
		testData := map[string]string{"key": "value"}
		fn, called, passedErr := createTestFallbackFunc(testData)
		ctx := context.Background()
		testError := errors.New("test error")

		handler := &fallbackHandler{fn: fn}

		// Act
		err := handler.handle(ctx, testError)

		// Assert
		assert.NoError(t, err, "Handle should not return error when fallback succeeds")
		assert.True(t, *called, "Fallback function should be called")
		assert.Equal(t, testError, *passedErr, "Original error should be passed to fallback")
		assert.True(t, handler.executed, "Handler should mark execution as completed")

		// Verify response structure
		require.NotNil(t, handler.resp, "Response should be created")
		assert.Equal(t, http.StatusOK, handler.resp.StatusCode)
		assert.Equal(t, "application/json", handler.resp.Header.Get("Content-Type"))

		// Verify response body
		respBody, err := io.ReadAll(handler.resp.Body)
		require.NoError(t, err)
		assert.Contains(t, string(respBody), `"key":"value"`, "Response body should contain our test data")
	})

	t.Run("failing fallback", func(t *testing.T) {
		// Arrange
		fn, called := createFailingFallbackFunc()
		ctx := context.Background()
		testError := errors.New("test error")

		handler := &fallbackHandler{fn: fn}

		// Act
		err := handler.handle(ctx, testError)

		// Assert
		assert.Error(t, err, "Handle should return error when fallback fails")
		assert.True(t, *called, "Fallback function should be called")
		assert.Equal(t, testError, err, "Original error should be returned when fallback fails")
		assert.False(t, handler.executed, "Handler should not mark execution as completed")
	})
}

func TestInterfaceToReadCloserWithLength(t *testing.T) {
	testCases := []struct {
		name        string
		data        interface{}
		expectError bool
	}{
		{
			name:        "map data",
			data:        map[string]string{"key": "value"},
			expectError: false,
		},
		{
			name:        "slice data",
			data:        []string{"value1", "value2"},
			expectError: false,
		},
		{
			name:        "struct data",
			data:        struct{ Name string }{"test"},
			expectError: false,
		},
		{
			name:        "channel data (not marshallable)",
			data:        make(chan int),
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Act
			rc, length, contentType, err := interfaceToReadCloserWithLength(tc.data)

			// Assert
			if tc.expectError {
				assert.Error(t, err, "Should return error for unmarshallable data")
				return
			}

			assert.NoError(t, err, "Should not return error for valid data")
			assert.NotNil(t, rc, "Should return a non-nil ReadCloser")
			assert.Greater(t, length, int64(0), "Content length should be greater than 0")
			assert.Equal(t, "application/json", contentType, "Content type should be application/json")

			// Read the body to verify it's valid JSON
			body, err := io.ReadAll(rc)
			require.NoError(t, err)
			assert.NotEmpty(t, body, "Body should not be empty")
		})
	}
}
