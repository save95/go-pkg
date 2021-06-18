package strutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReverse(t *testing.T) {
	yes := map[string]string{
		"abcdefg123": "321gfedcba",
		" abcdefg12": "21gfedcba ",
		"abcdefg 12": "21 gfedcba",
	}

	for s1, s2 := range yes {
		assert.Equal(t, Reverse(s1), s2)
	}
}
