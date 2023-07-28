package job

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/robfig/cron/v3"
	"github.com/save95/xlog"
)

type cronJobWrapper struct {
	maxRetry uint8 // 最大重试次数

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

	msg := fmt.Sprintf("[job] %s register", name)
	if w.log == nil {
		log.Print(msg)
	} else {
		w.log.Info(msg)
	}

	return &commandJob{
		jobName:  name,
		job:      job,
		args:     args,
		maxRetry: w.maxRetry,
		ctx:      w.ctx,
		log:      w.log,
	}
}
