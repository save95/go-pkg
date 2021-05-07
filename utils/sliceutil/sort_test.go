package sliceutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShuffle(t *testing.T) {
	source := []int{1, 2, 3, 4, 5}
	ns, ok := ToAny(source)
	assert.True(t, ok)

	Shuffle(ns)
	assert.NotEqual(t, source, ToInt(ns))
}
