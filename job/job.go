package job

import (
	"context"
	"errors"
	"log"

	"github.com/save95/xerror"
	"github.com/save95/xlog"
)

type commandJob struct {
	jobName string
	job     ICommandJob
	args    []string

	maxRetry uint8 // 最大重试次数

	failedSaver func(jobName string, in []string, err error)

	ctx context.Context
	log xlog.XLogger
}

func (j commandJob) Run() {
	j.logf("debug", "[job] %s run starting", j.jobName)
	defer j.logf("debug", "[job] %s run end", j.jobName)

	if err := j.job.Run(j.args...); nil != err {
		var xe xerror.XError
		if errors.As(err, &xe) {
			err = xe.Unwrap()
		}
		j.logf("error", "[job] %s run failed: %+v", j.jobName, err)

		if j.failedSaver != nil {
			j.failedSaver(j.jobName, j.args, err)
		}
	}
}

func (j commandJob) logf(level string, format string, args ...interface{}) {
	if nil == j.log {
		log.Printf(format, args...)
		return
	}

	switch level {
	case "debug":
		j.log.Debugf(format, args...)
	case "info":
		j.log.Infof(format, args...)
	case "warn":
		j.log.Warningf(format, args...)
	case "error", "err":
		j.log.Errorf(format, args...)
	}
}
