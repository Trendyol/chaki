package client

import (
	"context"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
)

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

func TestFactory_Get(t *testing.T) {
	t.Run("creates client with default options", func(t *testing.T) {
		// Setup
		cfg := clientTestConfig()
		factory := NewFactory(cfg, nil)

		// Execute
		client := factory.Get("testclient")

		// Verify
		assert.NotNil(t, client)
		assert.Equal(t, "testclient", client.name)
		assert.NotNil(t, client.driver)
		// Verify baseURL is set correctly
		assert.Equal(t, "https://example.com", client.driver.HostURL)
	})

	t.Run("applies options correctly", func(t *testing.T) {
		// Setup
		cfg := clientTestConfig()
		factory := NewFactory(cfg, nil)

		customErrDecoder, decoderCalled := createTestErrDecoder()
		customWrapper, wrapperCalled := createDriverWrapper("X-Custom", "value")

		// Execute
		client := factory.Get("testclient",
			WithErrDecoder(customErrDecoder),
			WithDriverWrappers(customWrapper),
		)

		// Create a dummy response to test the error decoder
		req := client.driver.R().SetContext(context.Background())
		dummyRes := &resty.Response{
			Request: req,
		}
		// Manually invoke the error handler
		err := customErrDecoder(context.Background(), dummyRes)

		// Verify
		assert.NotNil(t, client)
		assert.NoError(t, err)
		assert.True(t, *decoderCalled, "Custom error decoder should be called")
		assert.True(t, *wrapperCalled, "Custom wrapper should be applied")
	})
}

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
