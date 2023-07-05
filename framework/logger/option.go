package logger

import "github.com/sirupsen/logrus"

func WithFormat(format LogFormat) func(*logger) {
	return func(l *logger) {
		switch format {
		case LogFormatJson:
			l.formatter = &formatJson{}
		default:
			l.formatter = &formatText{}
		}
	}
}

func WithFormatter(formatter logrus.Formatter) func(*logger) {
	return func(l *logger) {
		l.formatter = formatter
	}
}
