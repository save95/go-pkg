package internal

import (
	"context"
	"fmt"
	"os"
	"sync/atomic"
	"time"

	"github.com/save95/go-pkg/framework/logger"
	"github.com/save95/xerror"
	"github.com/save95/xlog"

	"github.com/IBM/sarama"
)

type IConsumer interface {
	SetLogger(logger xlog.XLogger)
	SetPanicHandler(panicHandler func(interface{}))
	SetMaxRetry(maxRetry uint)
	SetFailedHandler(handler func(consumerGroup, topic string, msg []byte, err error))

	RegisterHandler(group string, topics []string, handler func(topic string, msg []byte) error) error
	Run()
	Close() error
}

type defaultConsumer struct {
	ctx   context.Context
	addrs []string

	maxRetry      uint
	panicHandler  func(interface{})
	logger        xlog.XLogger
	FailedHandler func(consumerGroup, topic string, msg []byte, err error)

	handlers []*groupHandlerParam

	running uint32 // 是否运行中 (running: 1 stopped: 0)

	cg     sarama.ConsumerGroup
	config *sarama.Config
}

func NewDefaultConsumer(ctx context.Context, addrs []string) IConsumer {
	hostname, _ := os.Hostname()

	conf := sarama.NewConfig()
	//conf.Version = sarama.V0_10_2_0
	//conf.Version = sarama.V2_1_0_0
	conf.Version = sarama.V3_6_0_0
	conf.ClientID = hostname

	//conf.Consumer.Return.Errors = false
	//conf.Consumer.IsolationLevel = sarama.ReadCommitted
	//conf.Consumer.Offsets.AutoCommit.Enable = false
	conf.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{
		sarama.NewBalanceStrategyRoundRobin(),
	}

	// 默认启动一个日志收集
	l := logger.NewConsoleLogger()
	sarama.Logger = newSaramaLogger(l)

	return &defaultConsumer{
		config:   conf,
		addrs:    addrs,
		ctx:      ctx,
		logger:   l,
		handlers: make([]*groupHandlerParam, 0),
	}
}

func (dcc *defaultConsumer) SetLogger(logger xlog.XLogger) {
	dcc.logger = logger
	sarama.Logger = newSaramaLogger(logger)
}

func (dcc *defaultConsumer) SetPanicHandler(panicHandler func(interface{})) {
	dcc.panicHandler = panicHandler
}

func (dcc *defaultConsumer) SetMaxRetry(maxRetry uint) {
	dcc.maxRetry = maxRetry
}

func (dcc *defaultConsumer) SetFailedHandler(handler func(consumerGroup, topic string, msg []byte, err error)) {
	dcc.FailedHandler = handler
}

type groupHandlerParam struct {
	ConsumerGroup sarama.ConsumerGroup
	Topics        []string
	Handler       sarama.ConsumerGroupHandler
}

func (dcc *defaultConsumer) RegisterHandler(group string, topics []string, handler func(topic string, msg []byte) error) error {
	for _, s := range topics {
		if len(s) == 0 {
			return xerror.New("topic must not be empty")
		}
	}
	fmt.Println(group, topics)

	// 启动 consumerGroup
	cg, err := sarama.NewConsumerGroup(dcc.addrs, group, dcc.config)
	if err != nil {
		return xerror.Wrap(err, "create consumer group client failed")
	}

	dcc.handlers = append(dcc.handlers, &groupHandlerParam{
		ConsumerGroup: cg,
		Topics:        topics,
		Handler: newConsumerGroupHandler(group, &groupHandlerConf{
			Logger:        dcc.logger,
			Handler:       handler,
			MaxRetry:      dcc.maxRetry,
			FailedHandler: dcc.FailedHandler,
		}),
	})
	return nil
}

func (dcc *defaultConsumer) Run() {
	if atomic.LoadUint32(&dcc.running) == 1 {
		dcc.logger.Warning("consume could not called while running")
		return
	}

	atomic.StoreUint32(&dcc.running, 1)

	for _, handler := range dcc.handlers {
		goWrap(handler, func(handler *groupHandlerParam) {
			dcc.handle(dcc.ctx, handler.ConsumerGroup, handler.Topics, handler.Handler)
		})
	}
}

func (dcc *defaultConsumer) handle(ctx context.Context, cg sarama.ConsumerGroup, topics []string, handler sarama.ConsumerGroupHandler) {
	dcc.logger.Debugf("topic consumer handle start, topics=%s", topics)
	for {
		if err := cg.Consume(ctx, topics, handler); err != nil {
			dcc.logger.Errorf("topic consume failed: [topics:%s] %s", topics, err)
			time.Sleep(3 * time.Second)
		}

		// check if context was cancelled, signaling that the consumer should stop
		if ctx.Err() != nil {
			dcc.logger.Errorf("topic consume context happen error, break: [topics:%s] %s", topics, ctx.Err())
			return
		}
	}
}

func (dcc *defaultConsumer) Close() error {
	atomic.StoreUint32(&dcc.running, 0)

	if nil != dcc.cg {
		return dcc.cg.Close()
	}

	return nil
}
