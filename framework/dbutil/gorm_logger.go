package dbutil

import (
	"time"

	"github.com/save95/xlog"
	"gorm.io/gorm/logger"
)

type dbWriter struct {
	log xlog.XLog
}

func (l *dbWriter) Printf(s string, vs ...interface{}) {
	l.log.Infof(s, vs...)
}

func newWriter(logger xlog.XLog) *dbWriter {
	return &dbWriter{log: logger}
}

func newLogger(l xlog.XLog) logger.Interface {
	return logger.New(
		newWriter(l),
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Info, // Log level
			IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
			Colorful:                  false,       // Disable color
		},
	)
}
