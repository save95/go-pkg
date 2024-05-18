package internal

import (
	"fmt"

	"github.com/save95/xlog"
)

type defaultFailedHandler struct {
	log xlog.XLogger
}

func newDefaultFailedHandler(log xlog.XLogger) *defaultFailedHandler {
	return &defaultFailedHandler{
		log: log,
	}
}

func (d defaultFailedHandler) Print(consumerGroup, topic string, msg []byte, err error) {
	tips := fmt.Sprintf("failed to consumer message to cg: %s, topic: %s, msg: %s, err: %+v", consumerGroup, topic, msg, err)
	if d.log == nil {
		fmt.Println(tips)
		return
	}

	d.log.Error(tips)
	return
}
