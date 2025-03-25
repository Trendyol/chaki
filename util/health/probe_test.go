package health

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_probe(t *testing.T) {
	errTest := errors.New("test error")
	tests := []struct {
		name        string
		probeType   string                          // "readiness" or "liveness"
		probeFunc   func(ctx context.Context) error // probe function, which can be nil
		expectedErr error
		shouldCall  bool // indicates whether the probe function is expected to be called
	}{
		// Readiness tests
		{
			name:        "readiness: valid case returns nil",
			probeType:   "readiness",
			probeFunc:   func(ctx context.Context) error { return nil },
			expectedErr: nil,
			shouldCall:  true,
		},
		{
			name:        "readiness: error case returns error",
			probeType:   "readiness",
			probeFunc:   func(ctx context.Context) error { return errTest },
			expectedErr: errTest,
			shouldCall:  true,
		},
		{
			name:        "readiness: nil function returns nil",
			probeType:   "readiness",
			probeFunc:   nil,
			expectedErr: nil,
			shouldCall:  false,
		},
		// Liveness tests
		{
			name:        "liveness: valid case returns nil",
			probeType:   "liveness",
			probeFunc:   func(ctx context.Context) error { return nil },
			expectedErr: nil,
			shouldCall:  true,
		},
		{
			name:        "liveness: error case returns error",
			probeType:   "liveness",
			probeFunc:   func(ctx context.Context) error { return errTest },
			expectedErr: errTest,
			shouldCall:  true,
		},
		{
			name:        "liveness: nil function returns nil",
			probeType:   "liveness",
			probeFunc:   nil,
			expectedErr: nil,
			shouldCall:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var called bool
			// If the probe function is not nil, wrap it to track if it was called.
			var wrappedFunc func(ctx context.Context) error
			if tt.probeFunc != nil {
				wrappedFunc = func(ctx context.Context) error {
					called = true
					return tt.probeFunc(ctx)
				}
			}

			// Create a probe instance with the appropriate probe function based on the probeType.
			var probe Probe
			switch tt.probeType {
			case "readiness":
				probe = NewProbe(nil, wrappedFunc)
			case "liveness":
				probe = NewProbe(wrappedFunc, nil)
			default:
				t.Fatalf("Unknown probe type: %s", tt.probeType)
			}

			// Call the appropriate probe function.
			var err error
			if tt.probeType == "readiness" {
				err = probe.Readiness(context.Background())
			} else {
				err = probe.Liveness(context.Background())
			}

			// Verify if the probe function was called as expected.
			if tt.shouldCall {
				assert.True(t, called, "probe function should have been called")
			} else {
				assert.False(t, called, "probe function should not have been called")
			}

			// Assert the expected error.
			if tt.expectedErr == nil {
				assert.NoError(t, err)
			} else {
				assert.Equal(t, tt.expectedErr, err)
			}
		})
	}
}
