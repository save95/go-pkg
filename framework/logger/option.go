package logger

import "github.com/sirupsen/logrus"

type Option func(*logger)

func WithFormat(format LogFormat) Option {
	return func(l *logger) {
		switch format {
		case LogFormatJson:
			l.formatter = &formatJson{}
		default:
			l.formatter = &formatText{}
		}
	}
}

func WithFormatter(formatter logrus.Formatter) Option {
	return func(l *logger) {
		l.formatter = formatter
	}
}
