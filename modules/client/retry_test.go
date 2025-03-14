package client

import (
	"testing"
	"time"

	"github.com/Trendyol/chaki/config"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestDelayType verifies that delay type string representations are correct
func TestDelayType(t *testing.T) {
	tests := []struct {
		name           string
		delayType      DelayType
		expectedString string
	}{
		{
			name:           "constant delay type",
			delayType:      ConstantDelay,
			expectedString: "constant",
		},
		{
			name:           "exponential delay type",
			delayType:      ExponentialDelay,
			expectedString: "exponential",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expectedString, string(tt.delayType))
		})
	}
}

// TestRetryPresets verifies that predefined retry configurations have the expected values
func TestRetryPresets(t *testing.T) {
	tests := []struct {
		name     string
		preset   *retryConfig
		validate func(*testing.T, *retryConfig)
	}{
		{
			name:   "default preset",
			preset: defaultRetryConfig,
			validate: func(t *testing.T, cfg *retryConfig) {
				assert.Equal(t, "default", cfg.Name)
				assert.Equal(t, 3, cfg.Count)
				assert.Equal(t, 100*time.Millisecond, cfg.Interval)
				assert.Equal(t, 5*time.Second, cfg.MaxDelay)
				assert.Equal(t, ConstantDelay, cfg.DelayType)
			},
		},
		{
			name:   "exponential preset",
			preset: exponentialRetryConfig,
			validate: func(t *testing.T, cfg *retryConfig) {
				assert.Equal(t, "exponential", cfg.Name)
				assert.Equal(t, 3, cfg.Count)
				assert.Equal(t, 100*time.Millisecond, cfg.Interval)
				assert.Equal(t, 5*time.Second, cfg.MaxDelay)
				assert.Equal(t, ExponentialDelay, cfg.DelayType)
			},
		},
		{
			name:   "aggressive preset",
			preset: aggressiveRetryConfig,
			validate: func(t *testing.T, cfg *retryConfig) {
				assert.Equal(t, "aggressive", cfg.Name)
				assert.Equal(t, 7, cfg.Count)
				assert.Equal(t, 50*time.Millisecond, cfg.Interval)
				assert.Equal(t, 2*time.Second, cfg.MaxDelay)
				assert.Equal(t, ConstantDelay, cfg.DelayType)
			},
		},
		{
			name:   "aggressive exponential preset",
			preset: aggressiveExponentialRetryConfig,
			validate: func(t *testing.T, cfg *retryConfig) {
				assert.Equal(t, "aggressiveExponential", cfg.Name)
				assert.Equal(t, 7, cfg.Count)
				assert.Equal(t, 50*time.Millisecond, cfg.Interval)
				assert.Equal(t, 2*time.Second, cfg.MaxDelay)
				assert.Equal(t, ExponentialDelay, cfg.DelayType)
			},
		},
		{
			name:   "relaxed preset",
			preset: relaxedRetryConfig,
			validate: func(t *testing.T, cfg *retryConfig) {
				assert.Equal(t, "relaxed", cfg.Name)
				assert.Equal(t, 2, cfg.Count)
				assert.Equal(t, 500*time.Millisecond, cfg.Interval)
				assert.Equal(t, 2*time.Second, cfg.MaxDelay)
				assert.Equal(t, ConstantDelay, cfg.DelayType)
			},
		},
		{
			name:   "relaxed exponential preset",
			preset: relaxedExponentialConfig,
			validate: func(t *testing.T, cfg *retryConfig) {
				assert.Equal(t, "relaxedExponential", cfg.Name)
				assert.Equal(t, 2, cfg.Count)
				assert.Equal(t, 500*time.Millisecond, cfg.Interval)
				assert.Equal(t, 2*time.Second, cfg.MaxDelay)
				assert.Equal(t, ExponentialDelay, cfg.DelayType)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Execute validation function
			tt.validate(t, tt.preset)
		})
	}
}

func TestGetRetryConfigs(t *testing.T) {
	t.Run("disabled retry", func(t *testing.T) {
		v := viper.New()
		cfg := config.NewConfig(v, nil)
		cfg.Set("retry.enabled", false)

		result := getRetryConfigs(cfg)
		assert.Nil(t, result)
	})

	t.Run("custom configuration", func(t *testing.T) {
		v := viper.New()
		cfg := config.NewConfig(v, nil)
		cfg.Set("retry.enabled", true)
		cfg.Set("retry.preset", "custom")
		cfg.Set("retry.name", "custom-retry")
		cfg.Set("retry.count", 5)
		cfg.Set("retry.interval", "200ms")
		cfg.Set("retry.maxDelay", "3s")
		cfg.Set("retry.delayType", "exponential")

		result := getRetryConfigs(cfg)
		require.NotNil(t, result)
		assert.Equal(t, "custom-retry", result.Name)
		assert.Equal(t, 5, result.Count)
		assert.Equal(t, 200*time.Millisecond, result.Interval)
		assert.Equal(t, 3*time.Second, result.MaxDelay)
		assert.Equal(t, ExponentialDelay, result.DelayType)
	})

	t.Run("preset configuration", func(t *testing.T) {
		v := viper.New()
		cfg := config.NewConfig(v, nil)
		cfg.Set("retry.enabled", true)
		cfg.Set("retry.preset", "default")

		result := getRetryConfigs(cfg)
		require.NotNil(t, result)
		assert.Equal(t, defaultRetryConfig, result)
	})

	t.Run("unknown preset", func(t *testing.T) {
		v := viper.New()
		cfg := config.NewConfig(v, nil)
		cfg.Set("retry.enabled", true)
		cfg.Set("retry.preset", "unknown")

		assert.Panics(t, func() {
			getRetryConfigs(cfg)
		})
	})
}

func TestSetDefaultRetryConfigs(t *testing.T) {
	// Arrange
	v := viper.New()
	cfg := config.NewConfig(v, nil)

	// Act
	setDefaultRetryConfigs(cfg)

	// Assert
	assert.False(t, cfg.GetBool("retry.enabled"))
	assert.Equal(t, "default", cfg.GetString("retry.preset"))
	assert.Equal(t, 3, cfg.GetInt("retry.count"))
	assert.Equal(t, "100ms", cfg.GetString("retry.interval"))
	assert.Equal(t, "5s", cfg.GetString("retry.maxDelay"))
	assert.Equal(t, "constant", cfg.GetString("retry.delayType"))
}

func TestInitRetryPresets(t *testing.T) {
	t.Run("predefined presets", func(t *testing.T) {
		// Arrange
		v := viper.New()
		cfg := config.NewConfig(v, nil)

		// Act
		initRetryPresets(cfg)

		// Assert - Check if all predefined presets are registered
		presets := []string{"default", "exponential", "aggressive", "aggressiveExponential", "relaxed", "relaxedExponential"}
		for _, preset := range presets {
			rc := retryPresetMap.Get(preset)
			assert.NotNil(t, rc, "Preset %s should be registered", preset)
		}
	})

	t.Run("custom presets", func(t *testing.T) {
		// Arrange
		v := viper.New()
		cfg := config.NewConfig(v, nil)
		customPreset := &retryConfig{
			Name:      "custom",
			Count:     5,
			Interval:  200 * time.Millisecond,
			MaxDelay:  3 * time.Second,
			DelayType: ExponentialDelay,
		}
		cfg.Set("client.retryPresets", []*retryConfig{customPreset})

		// Act
		initRetryPresets(cfg)

		// Assert
		rc := retryPresetMap.Get("custom")
		require.NotNil(t, rc)
		assert.Equal(t, customPreset, rc)
	})

	t.Run("invalid custom presets", func(t *testing.T) {
		// Arrange
		v := viper.New()
		cfg := config.NewConfig(v, nil)
		cfg.Set("client.retryPresets", "invalid") // Invalid type

		// Act & Assert
		assert.Panics(t, func() {
			initRetryPresets(cfg)
		})
	})
}

func TestDelayType_Values(t *testing.T) {
	tests := []struct {
		name      string
		delayType DelayType
		expected  string
	}{
		{
			name:      "constant delay",
			delayType: ConstantDelay,
			expected:  "constant",
		},
		{
			name:      "exponential delay",
			delayType: ExponentialDelay,
			expected:  "exponential",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expected, string(tc.delayType))
		})
	}
}
