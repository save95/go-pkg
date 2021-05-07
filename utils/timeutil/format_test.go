package timeutil

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestFormat(t *testing.T) {
	assert.Equal(t, "", Format(nil, "2006-01-02 15:04:05"))

	day := time.Date(2021, 2, 22, 22, 22, 22, 0, time.Local)
	assert.Equal(t, "2021-02-22 22:22:22", Format(&day, "2006-01-02 15:04:05"))
	assert.Equal(t, "2021/02/22", Format(&day, "2006/01/02"))
}
