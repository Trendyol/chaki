package chaki

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWithConfigDisabled(t *testing.T) {
	t.Run("get default disabled value when option provided", func(t *testing.T) {
		// Given
		var opts = []Option{WithConfigDisabled(true)}
		// When
		res := getOptions(opts...)
		// Then
		assert.Equal(t, true, res.configOptions.disabled)
	})
}

func TestWithConfigPath(t *testing.T) {
	t.Run("get default path when no option provided", func(t *testing.T) {
		// Given
		var opts = []Option{}
		// When
		res := getOptions(opts...)
		// Then
		assert.Equal(t, "resources/configs/application.yaml", res.configOptions.path)
	})

	t.Run("get default path when no option provided", func(t *testing.T) {
		// Given
		var want = "file/path.json"
		var opts = []Option{WithConfigPath(want)}
		// When
		res := getOptions(opts...)
		// Then
		assert.Equal(t, want, res.configOptions.path)
	})
}

func TestWithConfigReferncePath(t *testing.T) {
	t.Run("get default path when no option provided", func(t *testing.T) {
		// Given
		var opts = []Option{}
		// When
		res := getOptions(opts...)
		// Then
		assert.Equal(t, map[string]string{}, res.configOptions.referencePaths)
	})

	t.Run("get default path when no option provided", func(t *testing.T) {
		// Given
		var opts = []Option{
			WithConfigReferencePath("foo", "bar.json"),
			WithConfigReferencePath("zoo", "too.yaml"),
		}
		// When
		res := getOptions(opts...)
		// Then
		assert.Equal(t, "bar.json", res.configOptions.referencePaths["foo"])
		assert.Equal(t, "too.yaml", res.configOptions.referencePaths["zoo"])
	})
}
