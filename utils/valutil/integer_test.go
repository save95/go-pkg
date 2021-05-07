package valutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInt(t *testing.T) {
	assertInt(t, 1, 1.23, 1.56, "2.381", "2.6189", true, false)
	assertIntError(t, struct{}{})
}

func assertInt(t *testing.T, anys ...interface{}) {
	for i := range anys {
		v, err := Int(anys[i])
		assert.Nil(t, err)

		it := 1
		assert.IsType(t, it, v)
	}
}

func assertIntError(t *testing.T, anys ...interface{}) {
	for i := range anys {
		v, err := Int(anys[i])
		assert.Error(t, err)
		assert.Equal(t, v, 0)
	}
}
