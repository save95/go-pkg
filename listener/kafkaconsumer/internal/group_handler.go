package internal

import (
	"math"

	"github.com/IBM/sarama"
	"github.com/pkg/errors"
	"github.com/save95/xlog"
)

type groupHandler struct {
	consumerGroup string

	maxRetry      uint // 最大重试次数
	handler       func(topic string, msg []byte) error
	failedHandler func(consumerGroup, topic string, msg []byte, err error)

	logger xlog.XLogger
}

type groupHandlerConf struct {
	Logger xlog.XLogger

	Handler       func(topic string, msg []byte) error
	FailedHandler func(consumerGroup, topic string, msg []byte, err error)

	MaxRetry uint // 最大重试次数
}

func newConsumerGroupHandler(cg string, conf *groupHandlerConf) sarama.ConsumerGroupHandler {
	failedHandler := conf.FailedHandler
	// 如果没有定义错误处理函数，则默认使用日志打印
	if failedHandler == nil {
		failedHandler = newDefaultFailedHandler(conf.Logger).Print
	}

	return &groupHandler{
		consumerGroup: cg,
		handler:       conf.Handler,
		failedHandler: failedHandler,
		logger:        conf.Logger,
		maxRetry:      uint(math.Max(3, float64(conf.MaxRetry))),
	}
}

func (c *groupHandler) Setup(sarama.ConsumerGroupSession) error {
	//// 没有日志，则初始化一个默认日志处理
	//if nil == c.logger {
	//	l := console.NewLogger()
	//	l.SetLevel(int(xlog.DebugLevel))
	//	c.logger = l
	//	fmt.Println("not set logger, use default consoleLogger")
	//}

	return nil
}

func (c *groupHandler) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages().
func (c *groupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {

	for {
		select {
		case msg, ok := <-claim.Messages():
			if !ok {
				c.logger.Debug("message channel was closed")
				return nil
			}

			c.logger.Debugf(
				"message claimed: cg=%q, topic=%q, time=%v, partition=%d, offset=%d",
				c.consumerGroup, msg.Topic, msg.Timestamp, msg.Partition, msg.Offset,
			)

			if err := c.handler(msg.Topic, msg.Value); err != nil {
				//// todo 重试次数
				//if c.maxRetry != 0 && retryCount <= c.maxRetry {
				//	time.Sleep(time.Duration(math.Pow(3, float64(c.maxRetry))) * time.Second)
				//	retryCount++
				//	return errors.Wrapf(err, "event handle failed, wait to retry(%d)", retryCount)
				//}

				// 未定义失败处理函数，直接抛出错误，阻塞
				if c.failedHandler == nil {
					return errors.Wrapf(err, "event handle failed，data: %s", msg.Value)
				}

				c.failedHandler(c.consumerGroup, msg.Topic, msg.Value, err)
			}

			session.MarkMessage(msg, "")
		case <-session.Context().Done():
			// Should return when `session.Context()` is done.
			// If not, will raise `ErrRebalanceInProgress` or `read tcp <ip>:<port>: i/o timeout` when kafka rebalance. see:
			// https://github.com/IBM/sarama/issues/1192
			return nil
		}
	}
}
