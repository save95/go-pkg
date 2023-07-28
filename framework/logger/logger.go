package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/pkg/errors"
	"github.com/save95/xlog"
	"github.com/sirupsen/logrus"
)

type logger struct {
	path      string     // 日志存放路径
	category  string     // 日志分类
	stack     xlog.Stack // 日志存储方式
	level     xlog.Level // 日志等级
	engine    xlog.XLog  // 日志引擎
	formatter logrus.Formatter

	traceId string
}

func (l *logger) WithPreField(field xlog.XPreField) xlog.XLog {
	// todo
	return l
}

func (l *logger) WithField(key string, value interface{}, options ...interface{}) xlog.XLog {
	// todo
	return l
}

func (l *logger) WithFields(fields xlog.Fields, options ...interface{}) xlog.XLog {
	// todo
	return l
}

func (l *logger) SetFieldFormatter(f xlog.XFieldFormatter) {
	// todo
}

func NewDefaultLogger() xlog.XLogger {
	return NewLoggerWithTraceId("", defaultDir, defaultCategory, xlog.DailyStack)
}

func NewDefaultTraceLogger(traceId string) xlog.XLogger {
	return NewLoggerWithTraceId(traceId, defaultDir, defaultCategory, xlog.DailyStack)
}

func NewLogger(path, category string, stack xlog.Stack, opts ...Option) xlog.XLogger {
	return NewLoggerWithTraceId("", path, category, stack, opts...)
}

func NewLoggerWithTraceId(traceId, path, category string, stack xlog.Stack, opts ...Option) xlog.XLogger {
	logger := &logger{
		category: defaultCategory,
		traceId:  traceId,
	}

	if err := logger.setPath(path); nil != err {
		fmt.Printf("logger setPath failed: %s\n", err.Error())
	}
	if err := logger.setCategory(category); nil != err {
		fmt.Printf("logger setCategory failed: %s\n", err.Error())
	}
	if err := logger.setStack(stack); nil != err {
		fmt.Printf("logger setStack failed: %s\n", err.Error())
	}

	for _, opt := range opts {
		opt(logger)
	}

	if err := logger.setEngine(); nil != err {
		fmt.Printf("logger setEngine failed: %s\n", err.Error())
	}

	return logger
}

func (l *logger) setPath(path string) error {
	if path == "" {
		return errors.New("log path is empty")
	}

	l.path = strings.TrimRight(path, string(filepath.Separator))
	return nil
}

func (l *logger) getPath() string {
	if len(l.path) > 0 {
		return l.path
	}

	return defaultDir
}

func (l *logger) setCategory(category string) error {
	l.category = strings.Trim(category, string(filepath.Separator))

	// 如果 path 中最后已经含了 category，则不重复创建
	path := strings.TrimSuffix(l.path, l.category)
	l.path = path

	return nil
}

func (l *logger) getCategory() string {
	if len(l.category) > 0 {
		return l.category
	}

	return defaultCategory
}

func (l *logger) setStack(stack xlog.Stack) error {
	l.stack = stack

	return nil
}

func (l *logger) getFilenamePatten() string {
	filename := defaultFilenameFormat

	switch l.stack {
	case xlog.DailyStack:
		filename = "%Y-%m-%d.log"
	}

	return filename
}

func (l *logger) setEngine() error {
	if nil != l.engine {
		return nil
	}

	if l.formatter == nil {
		l.formatter = &formatText{}
	}

	// 初始化引擎
	eg := logrus.New()
	eg.SetFormatter(l.formatter)
	eg.SetLevel(logrus.InfoLevel)
	eg.SetOutput(os.Stdout)
	l.engine = eg

	// 创建目录
	path := fmt.Sprintf("%s/%s", l.getPath(), l.getCategory())
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err := os.MkdirAll(path, os.ModePerm)
		if err != nil {
			return errors.Wrap(err, "create log dir failed")
		}
	}

	// 打开文件
	fp := fmt.Sprintf("%s/%s", path, l.getFilenamePatten())
	rl, err := rotatelogs.New(fp)
	if nil != err {
		return errors.Wrap(err, "set logger rotate failed")
	}

	eg.SetOutput(rl)

	return nil
}

func (l *logger) GetStack() xlog.Stack {
	return l.stack
}

func (l *logger) logFormatMerge(format string) string {
	// 没有定义格式，则只显示字符串
	if len(format) == 0 {
		format = "%s"
	}

	// 没有 traceId 不合并
	if len(l.traceId) == 0 {
		return format
	}

	return fmt.Sprintf("[%s] %s", l.traceId, format)
}

func (l *logger) Info(args ...interface{}) {
	l.Infof("", args...)
}

func (l *logger) Infof(format string, args ...interface{}) {
	if nil == l.engine {
		_ = l.setEngine()
	}

	l.engine.Infof(l.logFormatMerge(format), args...)
}

func (l *logger) Debug(args ...interface{}) {
	l.Debugf("", args...)
}

func (l *logger) Debugf(format string, args ...interface{}) {
	if nil == l.engine {
		_ = l.setEngine()
	}

	l.engine.Debugf(l.logFormatMerge(format), args...)
}

func (l *logger) Warning(args ...interface{}) {
	l.Warningf("", args...)
}

func (l *logger) Warningf(format string, args ...interface{}) {
	if nil == l.engine {
		_ = l.setEngine()
	}

	l.engine.Warningf(l.logFormatMerge(format), args...)
}

func (l *logger) Error(args ...interface{}) {
	l.Errorf("", args...)
}

func (l *logger) Errorf(format string, args ...interface{}) {
	if nil == l.engine {
		_ = l.setEngine()
	}

	l.engine.Errorf(l.logFormatMerge(format), args...)
}

func (l *logger) SetLevel(level int) {
	l.level = xlog.Level(level)

	if nil == l.engine {
		return
	}

	if eg, ok := l.engine.(*logrus.Logger); ok {
		eg.SetLevel(logrus.Level(level))
	}
}

func (l *logger) SetLevelByString(level string) {
	lv := xlog.ParseLevel(level)
	l.SetLevel(int(lv))
}

func (l *logger) GetLevel() int {
	return int(l.level)
}

func (l *logger) SetStdPrint(b bool) bool {
	if nil != l.engine && b {
		if eg, ok := l.engine.(*logrus.Logger); ok {
			eg.AddHook(NewStdoutHook(eg.Formatter))
			return true
		}
	}

	return false
}
