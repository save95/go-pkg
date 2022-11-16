package listener

import (
	"context"
	"time"

	"github.com/save95/go-pkg/queue"
	"github.com/save95/xlog"
)

type redisConsumer struct {
	queueName string
	config    *queue.RedisQueueConfig

	ctx context.Context
	log xlog.XLogger

	fun func(val string) error
}

func (q *redisConsumer) Consume() error {
	queued := queue.NewSimpleRedis(q.config, q.queueName)

	for {
		str, err := queued.Pop(q.ctx)
		if nil != err {
			q.log.Warningf("get queue item failed: [%s]: %+v", q.queueName, err)
			continue
		}

		if len(str) == 0 {
			time.Sleep(5 * time.Second)
			continue
		}

		q.fun(str)
	}
}
