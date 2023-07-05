package logger

type LogFormat int8

const (
	LogFormatText LogFormat = iota
	LogFormatJson
)
