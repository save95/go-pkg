package logger

import (
	"fmt"
	"log"

	"github.com/save95/xlog"
)

type consoleLogger struct {
	stack xlog.Stack // 日志存储方式
	level xlog.Level // 日志等级

	traceId string
}

func (cl *consoleLogger) WithPreField(field xlog.XPreField) xlog.XLog {
	// todo
	return cl
}

func (cl *consoleLogger) WithField(key string, value interface{}, options ...interface{}) xlog.XLog {
	// todo
	return cl
}

func (cl *consoleLogger) WithFields(fields xlog.Fields, options ...interface{}) xlog.XLog {
	// todo
	return cl
}

func (cl *consoleLogger) SetFieldFormatter(f xlog.XFieldFormatter) {
	// todo
}

func NewConsoleLogger() xlog.XLogger {
	return NewConsoleLoggerWithTraceId("", xlog.DailyStack)
}

func NewConsoleTraceLogger(traceId string) xlog.XLogger {
	return NewConsoleLoggerWithTraceId(traceId, xlog.DailyStack)
}

func NewConsoleLoggerWith(stack xlog.Stack) xlog.XLogger {
	return NewConsoleLoggerWithTraceId("", stack)
}

func NewConsoleLoggerWithTraceId(traceId string, stack xlog.Stack) xlog.XLogger {
	return &consoleLogger{
		stack:   stack,
		traceId: traceId,
	}
}

func (cl *consoleLogger) GetStack() xlog.Stack {
	return cl.stack
}

func (cl *consoleLogger) logFormatMerge(format string) string {
	// 没有定义格式，则只显示字符串
	if len(format) == 0 {
		format = "%s"
	}

	// 没有 traceId 不合并
	if len(cl.traceId) == 0 {
		return format
	}

	return fmt.Sprintf("[%s] %s", cl.traceId, format)
}

func (cl *consoleLogger) Info(args ...interface{}) {
	log.Print(args...)
}

func (cl *consoleLogger) Infof(format string, args ...interface{}) {
	log.Printf(cl.logFormatMerge(format), args...)
}

func (cl *consoleLogger) Debug(args ...interface{}) {
	log.Print(args...)
}

func (cl *consoleLogger) Debugf(format string, args ...interface{}) {
	log.Printf(cl.logFormatMerge(format), args...)
}

func (cl *consoleLogger) Warning(args ...interface{}) {
	log.Print(args...)
}

func (cl *consoleLogger) Warningf(format string, args ...interface{}) {
	log.Printf(cl.logFormatMerge(format), args...)
}

func (cl *consoleLogger) Error(args ...interface{}) {
	log.Print(args...)
}

func (cl *consoleLogger) Errorf(format string, args ...interface{}) {
	log.Printf(cl.logFormatMerge(format), args...)
}

func (cl *consoleLogger) SetLevel(level int) {
	cl.level = xlog.Level(level)
}

func (cl *consoleLogger) SetLevelByString(level string) {
	lv := xlog.ParseLevel(level)
	cl.SetLevel(int(lv))
}

func (cl *consoleLogger) GetLevel() int {
	return int(cl.level)
}

func (cl *consoleLogger) SetStdPrint(b bool) bool {
	return true
}
