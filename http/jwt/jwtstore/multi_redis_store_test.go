package jwtstore

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMultiRedisStore(t *testing.T) {
	store := NewMultiRedisStore(_redisClient)

	// 多次登录
	for i, token := range _tokens {
		_ts = time.Now().Add(time.Minute * time.Duration((i+1)*10)).Unix()

		assert.Equal(t, token, _tokens[i])

		err := store.Save(_userID, token, _ts)
		assert.Nil(t, err)
	}

	// 判断 token
	for _, token := range _tokens {
		err := store.Check(_userID, token)
		assert.Nil(t, err)
	}

	err := store.Check(_userID, "not-in-tokens")
	assert.NotNil(t, err)

	// 删除 token
	err = store.Remove(_userID, "not-in-tokens")
	assert.Nil(t, err)

	err = store.Remove(_userID, _tokens[len(_tokens)-1])
	assert.Nil(t, err)

	// 清理 token
	err = store.Clean(_userID)
	assert.Nil(t, err)
}
