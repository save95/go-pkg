package singleconsumer

import (
	"context"
	"time"

	"github.com/save95/go-pkg/framework/logger"

	"github.com/save95/go-pkg/queue"
	"github.com/save95/xlog"
)

type redisConsumer struct {
	queueName string
	config    *queue.RedisQueueConfig

	log xlog.XLogger

	fun           func(val string) error
	failedHandler func(val string, err error)
}

func WithRedisConsumerLogger(log xlog.XLogger) func(*redisConsumer) {
	return func(c *redisConsumer) {
		c.log = log
	}
}

func WithRedisConsumerHandle(config *queue.RedisQueueConfig, queueName string, fun func(val string) error) func(*redisConsumer) {
	return func(c *redisConsumer) {
		c.config = config
		c.queueName = queueName
		c.fun = fun
	}
}

func WithRedisConsumerFailedHandler(handler func(val string, err error)) func(*redisConsumer) {
	return func(c *redisConsumer) {
		c.failedHandler = handler
	}
}

func NewRedisConsumer(opts ...func(*redisConsumer)) IConsumer {
	consumer := &redisConsumer{
		log: logger.NewConsoleLogger(),
		fun: func(val string) error {
			return nil
		},
	}
	for _, opt := range opts {
		opt(consumer)
	}
	return consumer
}

func (q *redisConsumer) Consume(ctx context.Context) error {
	queued := queue.NewSimpleRedis(q.config, q.queueName)

	for {
		str, err := queued.Pop(ctx)
		if nil != err {
			q.log.Warningf("get queue item failed: [%s]: %+v", q.queueName, err)
			continue
		}

		if len(str) == 0 {
			time.Sleep(5 * time.Second)
			continue
		}

		if err := q.fun(str); nil != err {
			if q.failedHandler != nil {
				q.failedHandler(str, err)
				continue
			}

			q.log.Warningf("handle queue item failed: [%s]: [%s] %+v", q.queueName, str, err)
		}
	}
}
