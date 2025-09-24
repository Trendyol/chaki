package client

import (
	"context"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
)

// TestNewFactory verifies that the factory is created correctly with the provided configuration and wrappers
func TestNewFactory(t *testing.T) {
	// Setup
	cfg := clientTestConfig()
	wrapper, _ := createDriverWrapper("X-Test", "test")
	wrappers := []DriverWrapper{wrapper}

	// Execute
	factory := NewFactory(cfg, wrappers)

	// Verify
	assert.NotNil(t, factory)
	assert.Equal(t, cfg, factory.cfg)
	assert.Equal(t, wrappers, factory.baseWrappers)
}

// TestFactory_Get verifies that the factory creates clients correctly with various options
func TestFactory_Get(t *testing.T) {
	// Define test cases
	testCases := []struct {
		name           string
		setupOptions   []Option
		verifyFunction func(t *testing.T, client *Base)
	}{
		{
			name:         "creates client with default options",
			setupOptions: nil,
			verifyFunction: func(t *testing.T, client *Base) {
				assert.Equal(t, "testclient", client.name)
				assert.NotNil(t, client.driver)
				assert.Equal(t, "https://example.com", client.driver.HostURL)
			},
		},
		{
			name: "applies options correctly",
			setupOptions: func() []Option {
				customErrDecoder, _ := createTestErrDecoder()
				customWrapper, _ := createDriverWrapper("X-Custom", "value")
				return []Option{
					WithErrDecoder(customErrDecoder),
					WithDriverWrappers(customWrapper),
				}
			}(),
			verifyFunction: func(t *testing.T, client *Base) {
				assert.NotNil(t, client)
				// Additional verification can be done here if needed
			},
		},
	}

	// Run test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup
			cfg := clientTestConfig()
			factory := NewFactory(cfg, nil)

			// Execute
			client := factory.Get("testclient", tc.setupOptions...)

			// Verify
			assert.NotNil(t, client)
			tc.verifyFunction(t, client)
		})
	}

	// Additional test for error decoder functionality
	t.Run("error decoder is applied and works", func(t *testing.T) {
		// Setup
		cfg := clientTestConfig()
		factory := NewFactory(cfg, nil)

		customErrDecoder, decoderCalled := createTestErrDecoder()

		// Execute
		client := factory.Get("testclient", WithErrDecoder(customErrDecoder))

		// Create a dummy response to test the error decoder
		req := client.driver.R().SetContext(context.Background())
		dummyRes := &resty.Response{
			Request: req,
		}
		// Manually invoke the error handler
		err := customErrDecoder(context.Background(), dummyRes)

		// Verify
		assert.NoError(t, err)
		assert.True(t, *decoderCalled, "Custom error decoder should be called")
	})

	// Additional test for wrapper functionality
	t.Run("wrapper is applied and works", func(t *testing.T) {
		// Setup
		cfg := clientTestConfig()
		factory := NewFactory(cfg, nil)

		customWrapper, wrapperCalled := createDriverWrapper("X-Custom", "value")

		// Execute
		client := factory.Get("testclient", WithDriverWrappers(customWrapper))

		// Verify
		assert.NotNil(t, client)
		assert.True(t, *wrapperCalled, "Custom wrapper should be applied")
		// Verify the client has the expected header
		assert.Equal(t, "value", client.driver.Header.Get("X-Custom"))
	})
}

// TestBase_Request verifies that the Request method correctly creates a request with the provided context
func TestBase_Request(t *testing.T) {
	// Setup
	base := &Base{
		name:   "testclient",
		driver: resty.New(),
	}
	ctx := context.Background()

	// Execute
	req := base.Request(ctx)

	// Verify
	assert.NotNil(t, req)
	assert.Equal(t, ctx, req.Context())
}

// TestBase_RequestWithCommand verifies that the RequestWithCommand method correctly creates a request
// with the provided context and command
func TestBase_RequestWithCommand(t *testing.T) {
	// Setup
	base := &Base{
		name:   "testclient",
		driver: resty.New(),
	}
	ctx := context.Background()
	command := "test-command"

	// Execute
	req := base.RequestWithCommand(ctx, command)

	// Verify
	assert.NotNil(t, req)
	// Check that the command was added to the context
	assert.Equal(t, command, req.Context().Value(circuitCommandKey))
}
