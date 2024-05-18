package singleconsumer

import (
	"context"

	"github.com/save95/xerror"

	"github.com/save95/go-pkg/framework/logger"
	"github.com/save95/xlog"
)

type server struct {
	ctx       context.Context
	logger    xlog.XLogger
	consumers []IConsumer
}

func New(ctx context.Context) *server {
	return &server{
		ctx:       ctx,
		logger:    logger.NewConsoleLogger(),
		consumers: make([]IConsumer, 0),
	}
}

func (s *server) Register(consumer IConsumer) {
	if consumer == nil {
		return
	}

	s.consumers = append(s.consumers, consumer)
}

func (s *server) CountConsumers() uint {
	return uint(len(s.consumers))
}

func (s *server) Start() error {
	if s.consumers == nil || len(s.consumers) == 0 {
		return xerror.New("no register consumers")
	}

	for _, consumer := range s.consumers {
		c := consumer
		go func() {
			_ = c.Consume(s.ctx)
		}()
	}

	return nil
}

func (s *server) Shutdown() error {
	//if s.c != nil {
	//	global.Log.Infof("listener server stop")
	//	s.c.Stop()
	//}

	return nil
}
