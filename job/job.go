package job

import (
	"context"
	"fmt"
	"log"

	"github.com/save95/xerror"
	"github.com/save95/xlog"
)

type commandJob struct {
	jobName string
	job     ICommandJob
	args    []string

	maxRetry uint8 // 最大重试次数

	ctx context.Context
	log xlog.XLogger
}

func (j commandJob) Run() {
	msg := fmt.Sprintf("[job] %s run starting", j.jobName)
	if nil == j.log {
		log.Print(msg)
	} else {
		j.log.Debug(msg)
	}
	defer func() {
		msg = fmt.Sprintf("[job] %s run end", j.jobName)
		if nil == j.log {
			log.Print(msg)
		} else {
			j.log.Debug(msg)
		}
	}()

	if err := j.job.Run(j.args...); nil != err {
		if xe, ok := err.(xerror.XError); ok {
			err = xe.Unwrap()
		}
		msg := fmt.Sprintf("[job] %s run failed: %+v", j.jobName, err)
		if nil == j.log {
			log.Print(msg)
		} else {
			j.log.Error(msg)
		}
	}
}
