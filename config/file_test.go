package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_parseFile(t *testing.T) {
	cases := []struct {
		Arg  string
		Dir  string
		File string
		Ext  string
	}{
		{"src/app/config.json", "src/app", "config", "json"},
	}

	for _, c := range cases {
		d, f, e := parseFile(c.Arg)
		assert.Equal(t, c.Ext, e)
		assert.Equal(t, c.File, f)
		assert.Equal(t, c.Dir, d)
	}
}
