package logger

import (
	"fmt"

	"github.com/sirupsen/logrus"
)

type stdoutHook struct {
	formatter logrus.Formatter
}

func (sh *stdoutHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (sh *stdoutHook) Fire(entry *logrus.Entry) error {
	bs, _ := sh.formatter.Format(entry)
	fmt.Println(string(bs))

	return nil
}

func NewStdoutHook(f logrus.Formatter) *stdoutHook {
	return &stdoutHook{formatter: f}
}
