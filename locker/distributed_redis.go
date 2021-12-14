package locker

import (
	"github.com/go-redis/redis/v8"
	"github.com/save95/go-utils/locker"
)

// NewDistributedRedisLock 创建分布式 redis 锁
func NewDistributedRedisLock(client *redis.Client) locker.ILocker {
	return locker.NewDistributedRedisLock(client)
}
