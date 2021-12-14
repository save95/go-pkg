package queue

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/save95/xerror"
)

type queue struct {
	name        string
	timeout     time.Duration
	redisClient *redis.Client
}

// RedisQueueConfig Redis 队列参数
type RedisQueueConfig struct {
	Addr     string
	Password string
	Timeout  time.Duration
}

// NewSimpleRedis 创建简单的 Redis 队列
func NewSimpleRedis(cnf *RedisQueueConfig, name string) IQueue {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     cnf.Addr,
		Password: cnf.Password,
		DB:       6,
	})

	timeout := 15 * time.Second
	if cnf.Timeout > 0 {
		timeout = cnf.Timeout
	}

	return &queue{
		name:        fmt.Sprintf("queue:%s", name),
		timeout:     timeout,
		redisClient: redisClient,
	}
}

func (q *queue) Push(ctx context.Context, value string) error {
	_, err := q.redisClient.LPush(ctx, q.name, value).Result()

	return err
}

func (q *queue) Pop(ctx context.Context) (string, error) {
	str, err := q.redisClient.BRPop(ctx, q.timeout, q.name).Result()
	if nil != err {
		if err == redis.Nil {
			return "", nil
		}
		return "", err
	}

	if len(str) != 2 {
		return "", xerror.New("queue value error")
	}

	return str[1], nil
}
