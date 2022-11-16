package listener

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/save95/go-pkg/httpsqs"
	"github.com/save95/xlog"
)

type httpSQSConsumer struct {
	handler httpsqs.IHandler

	ctx context.Context
	log xlog.XLogger

	retry int
}

func NewHttpSQSConsumer(handler httpsqs.IHandler) *httpSQSConsumer {
	return &httpSQSConsumer{
		handler: handler,

		ctx: context.Background(),
	}
}

func (s *httpSQSConsumer) WithContext(ctx context.Context) *httpSQSConsumer {
	if nil != ctx {
		s.ctx = ctx
	}

	return s
}

func (s *httpSQSConsumer) WithLog(log xlog.XLogger) *httpSQSConsumer {
	if nil != log {
		s.log = log
	}

	return s
}

func (s *httpSQSConsumer) Consume() error {
	msg := fmt.Sprintf("[httpsqs] %s consumer, start", s.handler.QueueName())
	client, err := s.handler.GetClient()
	if nil != err {
		msg = fmt.Sprintf("[httpsqs] %s consumer, start failed: %s", s.handler.QueueName(), err)
	}

	if s.log == nil {
		log.Print(msg)
	} else {
		s.log.Info(msg)
	}

	defer func() {
		msg := fmt.Sprintf("[httpsqs] %s consumer, end", s.handler.QueueName())
		if s.log == nil {
			log.Print(msg)
		} else {
			s.log.Info(msg)
		}
	}()

	for {
		if err := s.handler.OnBefore(s.ctx); nil != err {
			sleep := 2 << s.retry
			s.retry++
			msg := fmt.Sprintf("[httpsqs] %s onBefore failed, sleep %d minute: %+v", s.handler.QueueName(), sleep, err)
			if s.log == nil {
				log.Print(msg)
			} else {
				s.log.Errorf(msg)
			}
			time.Sleep(time.Duration(sleep) * time.Minute)
			continue
		}

		// 获得队列状态
		status, err := client.Status(s.ctx, s.handler.QueueName())
		if nil != err {
			msg := fmt.Sprintf("[httpsqs] %s get queue status failed: %+v", s.handler.QueueName(), err)
			if s.log == nil {
				log.Print(msg)
			} else {
				s.log.Errorf(msg)
			}
			//global.Log.Errorf("[httpsqs] %s get queue status failed: %+v", s.name, err)
			continue
		}

		// 消费完则跳过
		if status.Unread == 0 {
			time.Sleep(3 * time.Minute)
			continue
		}

		str, pos, err := client.Get(s.ctx, s.handler.QueueName())
		if nil != err {
			//global.Log.Errorf("[httpsqs] %s get queue item failed: %+v", s.name, err)
			msg := fmt.Sprintf("[httpsqs] %s get queue item failed: %+v", s.handler.QueueName(), err)
			if s.log == nil {
				log.Print(msg)
			} else {
				s.log.Errorf(msg)
			}
			continue
		}

		if len(str) == 0 {
			time.Sleep(5 * time.Second)
			continue
		}

		// 处理数据
		if err := s.handler.Handle(s.ctx, str, pos); nil != err {
			go func() {
				time.Sleep(3 * time.Second)

				s.handler.OnFailed(s.ctx, str, err)
			}()
		}
	}
}
