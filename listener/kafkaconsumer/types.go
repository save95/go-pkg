package kafkaconsumer

type IHandler func(topic string, msg []byte) error

type IRegister interface {
	Register(group string, handler IHandler, topic string, topics ...string)
	CountListener() uint
}

type IListener interface {
	Group() string
	Topics() []string
	Handler() IHandler
}

type listener struct {
	group   string
	topics  []string
	handler IHandler
}

func (l *listener) Group() string {
	return l.group
}

func (l *listener) Topics() []string {
	return l.topics
}

func (l *listener) Handler() IHandler {
	return l.handler
}

func NewListener(group string, handler IHandler, topic string, topics ...string) IListener {
	topics = append([]string{topic}, topics...)
	return &listener{
		group:   group,
		topics:  topics,
		handler: handler,
	}
}
