package listener

type wrapper struct {
}

func (w *wrapper) QueueConsumer(fun func(val string) error) IConsumer {
	return &queueConsumer{
		queueName:   "",
		queueConfig: nil,
		ctx:         nil,
		log:         nil,
		fun:         fun,
	}
}
