package internal

import (
	"github.com/save95/xlog"
)

type saramaLogger struct {
	logger xlog.XLogger
}

func newSaramaLogger(logger xlog.XLogger) *saramaLogger {
	return &saramaLogger{
		logger: logger,
	}
}

func (fb *saramaLogger) Println(v ...interface{}) {
	v = append([]interface{}{"[sarama] "}, v...)
	v = append(v, "\n")
	fb.logger.Debug(v...)
}

func (fb *saramaLogger) Printf(format string, v ...interface{}) {
	fb.logger.Debugf("[sarama] "+format, v...)
}

func (fb *saramaLogger) Print(v ...interface{}) {
	v = append([]interface{}{"[sarama] "}, v...)
	fb.logger.Debug(v...)
}
