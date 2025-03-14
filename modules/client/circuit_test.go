package client

import (
	"net/http"
	"testing"

	"github.com/Trendyol/chaki/config"
	"github.com/Trendyol/chaki/util/store"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestSetDefaultCircuitConfigs verifies that default circuit breaker configurations are set correctly
func TestSetDefaultCircuitConfigs(t *testing.T) {
	// Setup
	cfg := config.NewConfig(viper.New(), nil)

	// Execute
	setDefaultCircuitConfigs(cfg)

	// Verify default values are set correctly
	assert.Equal(t, false, cfg.GetBool("circuit.enabled"))
	assert.Equal(t, "default", cfg.GetString("circuit.preset"))
	assert.Equal(t, 5000, cfg.GetInt("circuit.timeout"))
	assert.Equal(t, 100, cfg.GetInt("circuit.maxConcurrentRequests"))
	assert.Equal(t, 20, cfg.GetInt("circuit.requestVolumeThreshold"))
	assert.Equal(t, 5000, cfg.GetInt("circuit.sleepWindow"))
	assert.Equal(t, 50, cfg.GetInt("circuit.errorPercentThreshold"))
}

// TestInitCircuitPresets verifies that circuit presets are initialized correctly
func TestInitCircuitPresets(t *testing.T) {
	// Reset circuit map before test
	circuitPresetMap = store.NewBucket[string, *circuitConfig](func(k string) *circuitConfig { return nil })

	t.Run("default presets initialization", func(t *testing.T) {
		// Setup
		cfg := config.NewConfig(viper.New(), nil)

		// Execute
		initCircuitPresets(cfg)

		// Verify default presets were registered
		assert.NotNil(t, circuitPresetMap.Get("default"))
		assert.NotNil(t, circuitPresetMap.Get("aggressive"))
		assert.NotNil(t, circuitPresetMap.Get("relaxed"))
	})

	t.Run("custom presets registration", func(t *testing.T) {
		// Setup
		customPreset := &circuitConfig{
			Name:                   "custom-test-preset",
			Enabled:                true,
			Timeout:                1000,
			MaxConcurrentRequests:  25,
			ErrorPercentThreshold:  30,
			RequestVolumeThreshold: 15,
			SleepWindow:            2000,
		}

		cfg := config.NewConfig(viper.New(), nil)
		cfg.Set("client.circuitPresets", []*circuitConfig{customPreset})

		// Reset and init again with new config
		circuitPresetMap = store.NewBucket[string, *circuitConfig](func(k string) *circuitConfig { return nil })

		// Execute
		initCircuitPresets(cfg)

		// Verify custom preset was registered
		preset := circuitPresetMap.Get("custom-test-preset")
		require.NotNil(t, preset)
		assert.Equal(t, customPreset.Timeout, preset.Timeout)
		assert.Equal(t, customPreset.MaxConcurrentRequests, preset.MaxConcurrentRequests)
		assert.Equal(t, customPreset.ErrorPercentThreshold, preset.ErrorPercentThreshold)
	})
}

// TestGetCircuitConfigs verifies that circuit configurations are retrieved correctly
func TestGetCircuitConfigs(t *testing.T) {
	defer resetHystrix()

	testCases := []struct {
		name           string
		configSetup    func(*config.Config)
		expectedResult *circuitConfig
	}{
		{
			name: "circuit disabled",
			configSetup: func(cfg *config.Config) {
				cfg.Set("circuit.enabled", false)
			},
			expectedResult: nil,
		},
		{
			name: "default preset",
			configSetup: func(cfg *config.Config) {
				cfg.Set("circuit.enabled", true)
				cfg.Set("circuit.preset", "default")
			},
			expectedResult: defaultCircuitConfig,
		},
		{
			name: "custom preset",
			configSetup: func(cfg *config.Config) {
				cfg.Set("circuit.enabled", true)
				cfg.Set("circuit.preset", "custom")
				cfg.Set("circuit.timeout", 1000)
				cfg.Set("circuit.maxConcurrentRequests", 50)
				cfg.Set("circuit.requestVolumeThreshold", 10)
				cfg.Set("circuit.sleepWindow", 2000)
				cfg.Set("circuit.errorPercentThreshold", 25)
			},
			expectedResult: &circuitConfig{
				Name:                   "custom",
				Enabled:                true,
				Timeout:                1000,
				MaxConcurrentRequests:  50,
				RequestVolumeThreshold: 10,
				SleepWindow:            2000,
				ErrorPercentThreshold:  25,
			},
		},
	}

	// Initialize circuit presets
	circuitPresetMap = store.NewBucket[string, *circuitConfig](func(k string) *circuitConfig { return nil })
	initCircuitPresets(config.NewConfig(viper.New(), nil))

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup
			cfg := config.NewConfig(viper.New(), nil)
			tc.configSetup(cfg)

			// Execute
			result := getCircuitConfigs(cfg)

			// Verify
			if tc.expectedResult == nil {
				assert.Nil(t, result, "Should return nil when circuit is disabled")
			} else {
				require.NotNil(t, result)
				assert.Equal(t, tc.expectedResult.Timeout, result.Timeout)
				assert.Equal(t, tc.expectedResult.MaxConcurrentRequests, result.MaxConcurrentRequests)
				assert.Equal(t, tc.expectedResult.RequestVolumeThreshold, result.RequestVolumeThreshold)
				assert.Equal(t, tc.expectedResult.SleepWindow, result.SleepWindow)
				assert.Equal(t, tc.expectedResult.ErrorPercentThreshold, result.ErrorPercentThreshold)
			}
		})
	}
}

func TestCircuitConfigToHystrixConfig(t *testing.T) {
	cc := &circuitConfig{
		Name:                   "test-config",
		Enabled:                true,
		Timeout:                1000,
		MaxConcurrentRequests:  50,
		ErrorPercentThreshold:  25,
		RequestVolumeThreshold: 10,
		SleepWindow:            2000,
	}

	hystrixCfg := cc.toHystrixConfig()

	assert.Equal(t, cc.Timeout, hystrixCfg.Timeout)
	assert.Equal(t, cc.MaxConcurrentRequests, hystrixCfg.MaxConcurrentRequests)
	assert.Equal(t, cc.ErrorPercentThreshold, hystrixCfg.ErrorPercentThreshold)
	assert.Equal(t, cc.RequestVolumeThreshold, hystrixCfg.RequestVolumeThreshold)
	assert.Equal(t, cc.SleepWindow, hystrixCfg.SleepWindow)
}

func TestShouldTreatStatusCodeAsFailure(t *testing.T) {
	tests := []struct {
		name              string
		config            circuitConfig
		statusCode        int
		expectedTreatment bool
	}{
		{
			name: "Default behavior with server error",
			config: circuitConfig{
				StatusCodeConfig: statusCodeConfig{},
			},
			statusCode:        http.StatusInternalServerError,
			expectedTreatment: true,
		},
		{
			name: "Default behavior with client error",
			config: circuitConfig{
				StatusCodeConfig: statusCodeConfig{},
			},
			statusCode:        http.StatusBadRequest,
			expectedTreatment: false,
		},
		{
			name: "Treat all error codes as failure",
			config: circuitConfig{
				StatusCodeConfig: statusCodeConfig{
					TreatAllErrorCodesAsFailure: true,
				},
			},
			statusCode:        http.StatusBadRequest,
			expectedTreatment: true,
		},
		{
			name: "Specific status code as failure",
			config: circuitConfig{
				StatusCodeConfig: statusCodeConfig{
					SpecificStatusCodes: []int{http.StatusBadRequest, http.StatusNotFound},
				},
			},
			statusCode:        http.StatusNotFound,
			expectedTreatment: true,
		},
		{
			name: "Ignore specific status code",
			config: circuitConfig{
				StatusCodeConfig: statusCodeConfig{
					IgnoreStatusCodes: []int{http.StatusInternalServerError},
				},
			},
			statusCode:        http.StatusInternalServerError,
			expectedTreatment: false,
		},
		{
			name: "Ignore takes precedence over specific",
			config: circuitConfig{
				StatusCodeConfig: statusCodeConfig{
					IgnoreStatusCodes:   []int{http.StatusNotFound},
					SpecificStatusCodes: []int{http.StatusNotFound},
				},
			},
			statusCode:        http.StatusNotFound,
			expectedTreatment: false,
		},
		{
			name: "Ignore takes precedence over treat all as failure",
			config: circuitConfig{
				StatusCodeConfig: statusCodeConfig{
					IgnoreStatusCodes:           []int{http.StatusBadRequest},
					TreatAllErrorCodesAsFailure: true,
				},
			},
			statusCode:        http.StatusBadRequest,
			expectedTreatment: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := tt.config.shouldTreatStatusCodeAsFailure(tt.statusCode)
			assert.Equal(t, tt.expectedTreatment, actual)
		})
	}
}

func TestDefaultCircuitConfigs(t *testing.T) {
	// Test default values
	assert.Equal(t, "default", defaultCircuitConfig.Name)
	assert.Equal(t, true, defaultCircuitConfig.Enabled)
	assert.Equal(t, 5000, defaultCircuitConfig.Timeout)
	assert.Equal(t, 100, defaultCircuitConfig.MaxConcurrentRequests)
	assert.Equal(t, 50, defaultCircuitConfig.ErrorPercentThreshold)
	assert.Equal(t, 20, defaultCircuitConfig.RequestVolumeThreshold)
	assert.Equal(t, 5000, defaultCircuitConfig.SleepWindow)

	// Test aggressive values
	assert.Equal(t, "aggressive", aggressiveCircuitConfig.Name)
	assert.Equal(t, 2000, aggressiveCircuitConfig.Timeout)
	assert.Equal(t, 50, aggressiveCircuitConfig.MaxConcurrentRequests)
	assert.Equal(t, 25, aggressiveCircuitConfig.ErrorPercentThreshold)
	assert.Equal(t, 10, aggressiveCircuitConfig.RequestVolumeThreshold)
	assert.Equal(t, 3000, aggressiveCircuitConfig.SleepWindow)

	// Test relaxed values
	assert.Equal(t, "relaxed", relaxedCircuitConfig.Name)
	assert.Equal(t, 10000, relaxedCircuitConfig.Timeout)
	assert.Equal(t, 200, relaxedCircuitConfig.MaxConcurrentRequests)
	assert.Equal(t, 75, relaxedCircuitConfig.ErrorPercentThreshold)
	assert.Equal(t, 40, relaxedCircuitConfig.RequestVolumeThreshold)
	assert.Equal(t, 7000, relaxedCircuitConfig.SleepWindow)
}

func TestCircuitPresetBucket(t *testing.T) {
	// Reset bucket before testing
	circuitPresetMap = store.NewBucket[string, *circuitConfig](func(k string) *circuitConfig { return nil })

	// Test default function returns nil for non-existent keys
	result := circuitPresetMap.Get("non-existent")
	assert.Nil(t, result)

	// Test setting and getting values
	testPreset := &circuitConfig{Name: "test-preset"}
	circuitPresetMap.Set("test-preset", testPreset)

	result = circuitPresetMap.Get("test-preset")
	assert.Equal(t, testPreset, result)

	// Test Has functionality
	assert.True(t, circuitPresetMap.Has("test-preset"))
	assert.False(t, circuitPresetMap.Has("non-existent"))

	// Test Remove functionality
	circuitPresetMap.Remove("test-preset")
	assert.False(t, circuitPresetMap.Has("test-preset"))
	assert.Nil(t, circuitPresetMap.Get("test-preset"))
}
