package nr

import (
	"testing"

	"github.com/Trendyol/chaki/config"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func Test_newApplication(t *testing.T) {
	t.Run("it should build disabled new relic instance when no config provided", func(t *testing.T) {
		// Given
		var (
			cfg     = config.NewConfig(viper.New(), nil)
			options []newrelic.ConfigOption
		)
		cfg.Set("newrelic", map[string]any{})

		// When
		nr, err := newApplication(cfg, options)

		// Then
		assert.NoError(t, err)
		nrc, ok := nr.Config()
		assert.True(t, ok)
		assert.False(t, nrc.Enabled)
	})

}

func Test_setDefaultConfigs(t *testing.T) {
	// Given
	cfg := config.NewConfig(viper.New(), nil)

	// When
	setDefaultConfigs(cfg)

	// Then
	assert.False(t, cfg.GetBool(keyEnabled))
	assert.True(t, cfg.GetBool(keyAppLogEnabled))
}
