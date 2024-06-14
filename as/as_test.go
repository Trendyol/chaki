package as

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// commmon types
type IFoo interface {
	Bar() string
}

// type-check
var _ IFoo = (*fooImpl)(nil)

type fooImpl struct {
	bar string
}

func (f *fooImpl) Bar() string {
	return f.bar
}

func Test_asser_Match(t *testing.T) {

	t.Run("it should match struct type correctly when correct constructor  provided", func(t *testing.T) {
		// Given
		asser := Struct[string]("testf")
		ctr := func() string {
			return "string"
		}

		// When
		res := asser.Match(ctr)

		// Then
		assert.Equal(t, true, res)
	})

	t.Run("it should not match struct type when wrong constructor provided", func(t *testing.T) {
		// Given
		asser := Struct[string]("testf")
		ctr := func() int {
			return 14
		}

		// When
		res := asser.Match(ctr)

		// Then
		assert.Equal(t, false, res)
	})

	t.Run("it should match interface when correct constructor provided", func(t *testing.T) {
		// Given
		asser := Interface[IFoo]("testf")
		ctr := func() *fooImpl {
			return &fooImpl{"bar"}
		}

		// When
		res := asser.Match(ctr)

		// Then
		assert.Equal(t, true, res)
	})

}
