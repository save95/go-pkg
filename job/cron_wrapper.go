package job

import (
	"context"
	"fmt"
	"strings"

	"github.com/robfig/cron/v3"
	"github.com/save95/xlog"
)

type cronJobWrapper struct {
	maxRetry uint8 // 最大重试次数

	failedSaver func(jobName string, in []string, err error) // 错误记录器

	ctx context.Context
	log xlog.XLogger
}

func NewCronJobWrapper(opts ...WrapperOption) IWrapper {
	w := &cronJobWrapper{}

	for _, opt := range opts {
		opt(w)
	}

	return w
}

func (w *cronJobWrapper) FromCommandJob(job ICommandJob, args ...string) cron.Job {
	name := strings.Trim(fmt.Sprintf("%T", job), "*")

	return &commandJob{
		jobName:     name,
		job:         job,
		args:        args,
		maxRetry:    w.maxRetry,
		failedSaver: w.failedSaver,
		ctx:         w.ctx,
		log:         w.log,
	}
}
