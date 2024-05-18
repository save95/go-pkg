package listener

import (
	"context"

	"github.com/save95/go-pkg/httpsqs"
	"github.com/save95/xlog"
)

type wrapper struct {
	ctx context.Context
	log xlog.XLogger
}

// Deprecated
func NewWrapper() *wrapper {
	return &wrapper{
		ctx: context.Background(),
	}
}

// Deprecated
func (w *wrapper) WithContext(ctx context.Context) *wrapper {
	if nil != ctx {
		w.ctx = ctx
	}
	return w
}

// Deprecated
func (w *wrapper) WithLog(log xlog.XLogger) *wrapper {
	if log != nil {
		w.log = log
	}
	return w
}

// Deprecated
func (w *wrapper) HTTPSQS(handler httpsqs.IHandler) IConsumer {
	return NewHttpSQSConsumer(handler).
		WithContext(w.ctx).
		WithLog(w.log)
}
