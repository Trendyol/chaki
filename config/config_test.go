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

func Test_ToStruct(t *testing.T) {
	t.Run("it should parse to struct succesfully", func(t *testing.T) {
		// Given
		type cfgStruct struct {
			Foo      string
			Bar      *string
			FooField string `mapstructure:"foo-field"`
			FooS     struct {
				FooFoo int
				FooBar bool
			}
		}

		v := viper.New()
		v.Set("p.foo", "foo_str")
		v.Set("p.bar", nil)
		v.Set("p.fooS.fooFoo", 12)
		v.Set("p.foo-field", "foo-kebab")

		cfg := NewConfig(v, nil)

		// When
		res, err := ToStruct[cfgStruct](cfg, "p")

		// Then
		assert.NoError(t, err)
		assert.Equal(t, "foo_str", res.Foo)
		assert.Nil(t, res.Bar)
		assert.Equal(t, 12, res.FooS.FooFoo)
		assert.Equal(t, false, res.FooS.FooBar)
		assert.Equal(t, "foo-kebab", res.FooField)

	})

	t.Run("it should parse reference values to struct succesfully", func(t *testing.T) {
		// Given
		type cfgStruct struct {
			Foo string
		}

		v := viper.New()
		v.Set("p.foo", "${bar:foo}")

		barref := viper.New()
		barref.Set("foo", "bar_foo")

		cfg := NewConfig(v, map[string]*viper.Viper{
			"bar": barref,
		})

		// When
		res, err := ToStruct[cfgStruct](cfg, "p")

		// Then
		assert.NoError(t, err)
		assert.Equal(t, "bar_foo", res.Foo)

	})
}
