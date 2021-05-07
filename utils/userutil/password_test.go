package userutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHash(t *testing.T) {
	hasher := NewHasher()
	assert.Equal(t, hasher.hash("123456"), "313233343536da39a3ee5e6b4b0d3255bfef95601890afd80709")

	pwd, err := hasher.Sum("123456")
	if nil != err {
		t.Error(err)
		return
	}
	//t.Logf("pwd: %s\n", pwd)
	assert.True(t, hasher.Check("123456", pwd))
	assert.False(t, hasher.Check("2222", pwd))
}
