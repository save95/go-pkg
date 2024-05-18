package singleconsumer

import (
	"context"
	"time"

	"github.com/save95/go-pkg/framework/logger"
	"github.com/save95/xerror"

	"github.com/save95/go-pkg/httpsqs"
	"github.com/save95/xlog"
)

type httpSQSConsumer struct {
	log xlog.XLogger

	handler httpsqs.IHandler
	retry   uint
}

func WithHttpSQSConsumerLogger(log xlog.XLogger) func(*httpSQSConsumer) {
	return func(s *httpSQSConsumer) {
		s.log = log
	}
}

func WithHttpSQSConsumerHandler(handler httpsqs.IHandler) func(*httpSQSConsumer) {
	return func(s *httpSQSConsumer) {
		s.handler = handler
	}
}

func WithHttpSQSConsumerRetry(count uint) func(*httpSQSConsumer) {
	return func(s *httpSQSConsumer) {
		s.retry = count
	}
}

func NewHttpSQSConsumer(opts ...func(consumer *httpSQSConsumer)) IConsumer {
	consumer := &httpSQSConsumer{
		log: logger.NewConsoleLogger(),
	}

	for _, opt := range opts {
		opt(consumer)
	}

	return consumer
}

func (s *httpSQSConsumer) Consume(ctx context.Context) error {
	if s.handler == nil {
		return xerror.New("no httpsqs handler")
	}

	client, err := s.handler.GetClient()
	if nil != err {
		return xerror.Wrap(err, "get httpsqs client failed")
	}

	s.log.Debugf("[httpsqs] %s consumer, start", s.handler.QueueName())
	defer func() {
		s.log.Debugf("[httpsqs] %s consumer, end", s.handler.QueueName())
	}()

	for {
		if err := s.handler.OnBefore(ctx); nil != err {
			sleep := 2 << s.retry
			s.retry++
			s.log.Errorf("[httpsqs] %s onBefore failed, sleep %d minute: %+v", s.handler.QueueName(), sleep, err)
			time.Sleep(time.Duration(sleep) * time.Minute)
			continue
		}

		// 获得队列状态
		status, err := client.Status(ctx, s.handler.QueueName())
		if nil != err {
			s.log.Errorf("[httpsqs] %s get queue status failed: %+v", s.handler.QueueName(), err)
			continue
		}

		// 消费完则跳过
		if status.Unread == 0 {
			time.Sleep(3 * time.Minute)
			continue
		}

		str, pos, err := client.Get(ctx, s.handler.QueueName())
		if nil != err {
			s.log.Errorf("[httpsqs] %s get queue item failed: %+v", s.handler.QueueName(), err)
			continue
		}

		if len(str) == 0 {
			time.Sleep(5 * time.Second)
			continue
		}

		// 处理数据
		if err := s.handler.Handle(ctx, str, pos); nil != err {
			go func() {
				time.Sleep(3 * time.Second)

				s.handler.OnFailed(ctx, str, err)
			}()
		}
	}
}
