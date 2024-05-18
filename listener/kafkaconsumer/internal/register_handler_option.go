package internal

func WithHandler(handler func(topic string, msg []byte) error) func(*groupHandler) {
	return func(h *groupHandler) {
		h.handler = handler
	}
}

func WithFailedHandler(handler func(topic string, msg []byte) error) func(*groupHandler) {
	return func(h *groupHandler) {
		h.handler = handler
	}
}
