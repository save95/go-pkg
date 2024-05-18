package singleconsumer

import (
	"context"
)

// IConsumer 消费者约定
type IConsumer interface {
	Consume(ctx context.Context) error
}

type IRegister interface {
	Register(consumer IConsumer)
	CountConsumers() uint
}
