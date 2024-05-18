package kafkaconsumer

import (
	"context"

	"github.com/save95/xerror"

	"github.com/save95/go-pkg/framework/logger"

	"github.com/save95/go-pkg/listener/kafkaconsumer/internal"

	"github.com/save95/xlog"
)

type server struct {
	ctx    context.Context
	logger xlog.XLogger

	addrs     []string
	listeners []IListener

	maxRetry      uint
	panicHandler  func(interface{})
	failedHandler func(consumerGroup, topic string, msg []byte, err error)

	consumer internal.IConsumer
}

func New(addrs []string, opts ...Option) *server {
	svr := &server{
		addrs:     addrs,
		listeners: make([]IListener, 0),
		logger:    logger.NewConsoleLogger(),
	}

	for _, opt := range opts {
		opt(svr)
	}

	return svr
}

func (s *server) Register(group string, handler IHandler, topic string, topics ...string) {
	topics = append([]string{topic}, topics...)
	s.listeners = append(s.listeners, &listener{
		group:   group,
		topics:  topics,
		handler: handler,
	})
}

func (s *server) CountListener() uint {
	return uint(len(s.listeners))
}

func (s *server) Start() error {
	if s.listeners == nil || len(s.listeners) == 0 {
		return xerror.New("no register listener")
	}

	s.consumer = internal.NewDefaultConsumer(s.ctx, s.addrs)
	s.consumer.SetLogger(s.logger)
	s.consumer.SetMaxRetry(s.maxRetry)
	s.consumer.SetPanicHandler(s.panicHandler)
	s.consumer.SetFailedHandler(s.failedHandler)

	for _, item := range s.listeners {
		if err := s.consumer.RegisterHandler(item.Group(), item.Topics(), item.Handler()); nil != err {
			return err
		}
	}

	s.consumer.Run()
	return nil
}

func (s *server) Shutdown() error {
	if s.consumer != nil {
		_ = s.consumer.Close()
	}

	return nil
}
