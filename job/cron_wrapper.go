package job

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/robfig/cron/v3"
	"github.com/save95/xlog"
)

type cronJob struct {
	jobName  string
	job      IJob
	maxRetry uint8 // 最大重试次数

	ctx context.Context
	log xlog.XLogger
}

func (j cronJob) Run() {
	msg := fmt.Sprintf("[job] %s run starting", j.jobName)
	if nil == j.log {
		log.Print(msg)
	} else {
		j.log.Debug(msg)
	}

	if err := j.job.Run(); nil != err {
		msg := fmt.Sprintf("[job] %s run failed: %+v", j.jobName, err)
		if nil == j.log {
			log.Print(msg)
		} else {
			j.log.Error(msg)
		}
	}

	msg = fmt.Sprintf("[job] %s run end", j.jobName)
	if nil == j.log {
		log.Print(msg)
	} else {
		j.log.Debug(msg)
	}
}

// CronWrapper job 包装方法
// Deprecated
func CronWrapper(j IJob) cron.Job {
	return CronRetryWrapper(j, 0)
}

// CronRetryWrapper job 重试包装方法
// Deprecated
func CronRetryWrapper(j IJob, retry uint8) cron.Job {
	return &cronJob{
		job:      j,
		maxRetry: retry,
	}
}

type wrapper struct {
	maxRetry uint8 // 最大重试次数
	retry    uint8 // 重试次数

	ctx context.Context
	log xlog.XLogger
}

func NewWrapper() *wrapper {
	return &wrapper{}
}

func (w *wrapper) WithContext(ctx context.Context) *wrapper {
	w.ctx = ctx
	return w
}

func (w *wrapper) WithLog(log xlog.XLogger) *wrapper {
	w.log = log
	return w
}

func (w *wrapper) WithMaxRetry(retry uint8) *wrapper {
	w.maxRetry = retry
	return w
}

func (w *wrapper) Cron(job IJob) cron.Job {
	name := strings.Trim(fmt.Sprintf("%T", job), "*")

	msg := fmt.Sprintf("[job] %s register", name)
	if w.log == nil {
		log.Print(msg)
	} else {
		w.log.Info(msg)
	}

	return &cronJob{
		jobName:  name,
		job:      job,
		maxRetry: w.maxRetry,
		ctx:      w.ctx,
		log:      w.log,
	}
}
