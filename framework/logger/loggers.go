package logger

import (
	"sync"

	"github.com/save95/xlog"
)

type loggers struct {
	path     string
	stack    xlog.Stack
	level    xlog.Level
	stdPrint bool

	engines sync.Map // 存储所有 logger
}

func NewLoggers(path string, categories []string, stack xlog.Stack) *loggers {
	if len(path) == 0 {
		path = defaultDir
	}

	loggers := &loggers{
		path:  path,
		stack: stack,
	}

	if len(categories) == 0 {
		categories = append(categories, defaultCategory)
	}

	for _, category := range categories {
		loggers.engines.Store(category, NewLogger(path, category, stack))
	}

	return loggers
}

func (l *loggers) SetLevel(level int) {
	l.level = xlog.Level(level)
}

func (l *loggers) SetLevelByString(level string) {
	lv := xlog.ParseLevel(level)

	l.SetLevel(int(lv))
}

func (l *loggers) GetLevel() int {
	return int(l.level)
}

func (l *loggers) SetStdPrint(b bool) {
	l.stdPrint = b
}

func (l *loggers) GetLogger(category string) xlog.XLogger {
	var rl xlog.XLogger

	if lg, ok := l.engines.Load(category); ok {
		rl = lg.(*logger)
	} else {
		// 没有，则立即创建一个
		if len(l.path) == 0 {
			l.path = defaultDir
		}

		rl = NewLogger(l.path, category, l.stack)
		l.engines.Store(category, rl)
	}

	rl.SetLevel(int(l.level))
	rl.SetStdPrint(l.stdPrint)

	return rl
}
