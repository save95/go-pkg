package kafkapusher

import "context"

type IPusher interface {
	// Produce 生产消息
	Produce(ctx context.Context, topic string, message []byte) error
	// Produces 批量生产消息
	Produces(ctx context.Context, topic string, message ...[]byte) error
	// ProduceWithSequence 保持顺序性的生产消息
	// sequenceKey 顺序性特定KEY。如 kafka 为 partition key
	ProduceWithSequence(ctx context.Context, topic, sequenceKey string, messages ...[]byte) error
}
