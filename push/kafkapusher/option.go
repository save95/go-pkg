package kafkapusher

import (
	"time"

	"github.com/save95/xlog"
)

//func WithBrokers(brokers []string) func(*producer) {
//	return func(p *producer) {
//		p.brokers = brokers
//	}
//}

func WithTimeout(timeout time.Duration) func(*producer) {
	return func(p *producer) {
		p.timeout = timeout
	}
}

func WithLogger(logger xlog.XLogger) func(*producer) {
	return func(p *producer) {
		p.logger = logger
	}
}
