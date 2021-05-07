package valutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBool(t *testing.T) {
	assertTrue(t, "true", "TRue", "yes", "1", 1, 2, 3.14)
	assertFalse(t, "false", "FALSE", "no", "", "0", 0, -1, -1.23)
	assertBooleanError(t, "hello", "world", struct{}{}, []int{})
}

func assertTrue(t *testing.T, anys ...interface{}) {
	for i := range anys {
		v, err := Bool(anys[i])
		assert.Nil(t, err)
		assert.True(t, v)
	}
}

func assertFalse(t *testing.T, anys ...interface{}) {
	for i := range anys {
		v, err := Bool(anys[i])
		assert.Nil(t, err)
		assert.False(t, v)
	}
}

func assertBooleanError(t *testing.T, anys ...interface{}) {
	for i := range anys {
		v, err := Bool(anys[i])
		assert.Error(t, err)
		assert.False(t, v)
	}
}
