package dbutil

import (
	"github.com/jinzhu/gorm"
	"github.com/save95/xlog"
)

type dbLog struct {
	log xlog.XLog
}

func (l *dbLog) Print(v ...interface{}) {
	l.log.Info(gorm.LogFormatter(v...)...)
}

func convertLogger(logger xlog.XLog) *dbLog {
	return &dbLog{log: logger}
}
