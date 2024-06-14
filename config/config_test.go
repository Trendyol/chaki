package config

import (
	"os"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func Test_newEnvViper(t *testing.T) {
	t.Run("it should get config from env", func(t *testing.T) {
		// Given
		var (
			key   = "ENV_FOO"
			value = "ENV_BAR"
		)
		os.Setenv(key, value)
		v := newEnvViper()

		// When
		got := v.GetString(key)

		// Then
		assert.Equal(t, value, got)

	})

	t.Run("it should get config from env when set after init", func(t *testing.T) {
		// Given
		var (
			key   = "ENV_FOO"
			value = "ENV_BAR"
		)
		v := newEnvViper()
		os.Setenv(key, value)

		// When
		got := v.GetString(key)

		// Then
		assert.Equal(t, value, got)

	})
}

func Test_NewConfig(t *testing.T) {
	// Given
	// env variables
	os.Setenv("FOO_ENV", "bar_env")
	// reference variable
	refv := viper.New()
	refv.Set("foo", "bar_ref")
	// config
	cfgv := viper.New()
	cfgv.Set("foo", "bar")
	cfgv.Set("env_foo", "${env:FOO_ENV}")
	cfgv.Set("ref_foo", "${ref:foo}")

	// When
	cfg := NewConfig(cfgv, map[string]*viper.Viper{
		"ref": refv,
	})

	// Then
	assert.Equal(t, "bar_env", cfg.GetString("env_foo"))
	assert.Equal(t, "bar_ref", cfg.GetString("ref_foo"))
	assert.Equal(t, "bar", cfg.GetString("foo"))
}
