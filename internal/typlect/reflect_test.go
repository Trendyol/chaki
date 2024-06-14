package typlect

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsInterface(t *testing.T) {
	type a struct {
		Foo string
	}

	type b interface {
		Foo() string
	}

	type ft func() string

	assert.False(t, IsInterface[a]())
	assert.True(t, IsInterface[b]())
	assert.False(t, IsInterface[ft]())
}
