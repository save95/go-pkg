package listener

// IConsumer 消费者约定
type IConsumer interface {
	Consume() error
}
