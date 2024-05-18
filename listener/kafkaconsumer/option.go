package kafkaconsumer

import (
	"context"

	"github.com/save95/xlog"
)

type Option func(*server)

func WithContext(ctx context.Context) Option {
	return func(server *server) {
		server.ctx = ctx
	}
}

func WithLogger(logger xlog.XLogger) Option {
	return func(server *server) {
		server.logger = logger
	}
}

//func WithMaxRetry(retry int) Option {
//	return func(server *server) {
//		server.maxRetry = uint(retry)
//	}
//}

func WithPanicHandler(handler func(interface{})) Option {
	return func(server *server) {
		server.panicHandler = handler
	}
}

func WithConsumeGroupFailedHandler(handler func(consumerGroup, topic string, msg []byte, err error)) Option {
	return func(server *server) {
		server.failedHandler = handler
	}
}

func WithListeners(listeners []IListener) Option {
	return func(server *server) {
		server.listeners = listeners
	}
}

func WithListener(listener IListener) Option {
	return func(server *server) {
		if server.listeners == nil {
			server.listeners = make([]IListener, 0)
		}
		server.listeners = append(server.listeners, listener)
	}
}
