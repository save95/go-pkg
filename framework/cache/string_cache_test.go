package cache

import (
	"testing"

	"github.com/eko/gocache/v2/cache"
	"github.com/eko/gocache/v2/store"
	"github.com/go-redis/redis/v8"
)

var (
	cacheManager *cache.Cache
)

func init() {
	cacheManager = cache.New(store.NewRedis(redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "",
		DB:       4,
	}), nil))
}

func TestNewString(t *testing.T) {
	// todo
}
