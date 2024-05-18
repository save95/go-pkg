package kafkapusher

import (
	"context"
	"os"
	"sync"
	"time"

	"github.com/save95/xerror"

	"github.com/save95/go-pkg/framework/logger"

	"github.com/save95/xlog"

	"github.com/IBM/sarama"
)

var once sync.Once

type producer struct {
	brokers []string
	timeout time.Duration

	logger xlog.XLogger

	saramaConfig *sarama.Config
	producer     sarama.SyncProducer
}

func New(brokers []string, opts ...func(*producer)) *producer {
	p := &producer{
		brokers: brokers,
		timeout: time.Second * 5,
		logger:  logger.NewConsoleLogger(),
	}

	for _, opt := range opts {
		opt(p)
	}

	return p
}

func (b *producer) Produce(ctx context.Context, topic string, message []byte) error {
	return b.ProduceWithSequence(ctx, topic, "", message)
}

func (b *producer) Produces(ctx context.Context, topic string, message ...[]byte) error {
	return b.ProduceWithSequence(ctx, topic, "", message...)
}

// ProduceWithSequence 保持顺序性的生产消息
// sequenceKey 顺序性特定KEY。如 kafka 为 partition key
func (b *producer) ProduceWithSequence(ctx context.Context, topic, sequenceKey string, messages ...[]byte) error {
	if len(messages) == 0 {
		return xerror.New("no message")
	}

	msgs := make([]*sarama.ProducerMessage, 0, len(messages))
	for _, message := range messages {
		msgs = append(msgs, &sarama.ProducerMessage{
			Topic: topic,
			Key:   sarama.StringEncoder(sequenceKey),
			Value: sarama.ByteEncoder(message),
			//Headers:   nil,
			//Metadata:  nil,
			//Offset:    0,
			//Partition: 0,
			//Timestamp: time.Time{},
		})
	}

	return b.sends(ctx, msgs...)
}

func (b *producer) getSaramaConfig() *sarama.Config {
	if b.saramaConfig != nil {
		return b.saramaConfig
	}

	hostname, _ := os.Hostname()

	conf := sarama.NewConfig()

	// client common config
	conf.Version = sarama.V3_6_0_0
	conf.ClientID = hostname

	// producer config
	conf.Producer.RequiredAcks = sarama.WaitForAll
	conf.Producer.Return.Successes = true                  //接收 producer 的 success
	conf.Producer.Return.Errors = true                     // 接收 producer 的 error
	conf.Producer.Compression = sarama.CompressionZSTD     // 压缩方式，如果kafka版本大于1.2，推荐使用zstd压缩
	conf.Producer.Flush.Messages = 10                      // 缓存条数
	conf.Producer.Flush.Frequency = 500 * time.Millisecond // 缓存时间
	conf.Producer.Partitioner = sarama.NewRoundRobinPartitioner
	conf.Producer.Transaction.Retry.Backoff = 10

	// 日志收集
	sarama.Logger = newSaramaLogger(b.logger)

	return conf
}

func (b *producer) getProducer() sarama.SyncProducer {
	once.Do(func() {
		if b.producer == nil {
			if p, err := sarama.NewSyncProducer(b.brokers, b.getSaramaConfig()); err != nil {
				b.logger.Errorf("producer create failed: %v", err)
			} else {
				b.producer = p
			}
		}
	})

	return b.producer
}

func (b *producer) sends(_ context.Context, msgs ...*sarama.ProducerMessage) error {
	p := b.getProducer()
	if nil == p {
		return xerror.New("no producer found")
	}

	return p.SendMessages(msgs)
}
