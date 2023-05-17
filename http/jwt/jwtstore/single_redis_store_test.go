package jwtstore

import (
	"testing"

	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
)

var _redisClient *redis.Client

var (
	_userID = uint(100008)
	_tokens = []string{
		"tIqsOkAqXCum1AhiCTAMB4GqNmduU63l-1",
		"tIqsOkAqXCum1AhiCTAMB4GqNmduU63l-2",
		"tIqsOkAqXCum1AhiCTAMB4GqNmduU63l-3",
		"tIqsOkAqXCum1AhiCTAMB4GqNmduU63l-4",
	}
	_ts = int64(86400)
)

func init() {
	_redisClient = redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6379",
		DB:   2,
	})
}

func TestSingleRedisStore(t *testing.T) {
	store := NewSingleRedisStore(_redisClient)

	// 多次登录
	for _, token := range _tokens {
		err := store.Save(_userID, token, _ts)
		assert.Nil(t, err)
	}

	// 判断 token
	for i, token := range _tokens {
		err := store.Check(_userID, token)
		if i != len(_tokens)-1 {
			assert.NotNil(t, err)
		} else {
			assert.Nil(t, err)
		}
	}

	// 删除 token
	err := store.Remove(_userID, "error-token")
	assert.NotNil(t, err)

	err = store.Remove(_userID, _tokens[len(_tokens)-1])
	assert.Nil(t, err)

	// 清理 token
	err = store.Clean(_userID)
	assert.Nil(t, err)
}
