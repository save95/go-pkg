package sliceutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIs(t *testing.T) {
	a := []int{1, 2}
	assert.True(t, Is(a))

	b := []interface{}{1, "a"}
	assert.True(t, Is(b))

	c := 1
	assert.False(t, Is(c))

	d := [2]int{1, 2}
	assert.False(t, Is(d))
}
