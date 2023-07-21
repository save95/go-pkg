package jwt

import (
	"testing"

	"github.com/save95/go-pkg/http/types"
	"github.com/stretchr/testify/assert"
)

func TestToken(t *testing.T) {
	tk := NewToken(types.User{
		ID:      1,
		Account: "account",
		Name:    "My Name",
		Roles:   nil,
		IP:      "127.0.0.1",
		Extend: map[string]string{
			"a": "124",
		},
	})

	t.Log(tk.ToString())

	tks := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2ODk4MzQ4MDIsImlhdCI6MTY4OTc0ODQwMiwiaXNzIjoiZ28tcGtnIiwiYWNjb3VudCI6ImFjY291bnQiLCJ1aWQiOjEsIm5hbWUiOiJNeSBOYW1lIiwicm9sZXMiOm51bGwsImlwIjoiMTI3LjAuMC4xIiwiZXh0ZW5kIjp7ImEiOiIxMjQifX0.w0h9Kbqt4xY9fCvPZf5kg-WiFPbPCZ9oO99ybYjtRCg"
	c, err := parseToken(tks, jwtSecret)
	assert.Nil(t, err)

	ntks, err := newTokenWith(c).ToString()
	assert.Nil(t, err)
	assert.NotEqual(t, tks, ntks)
}
