package job

import (
	"log"

	"github.com/robfig/cron/v3"
)

type cronJob struct {
	job      IJob
	maxRetry uint8 // 最大重试次数
	retry    uint8 // 重试次数
}

func (cjw cronJob) Run() {
	err := cjw.job.Run()
	log.Printf("job failed: %+v\n", err)
}

// CronWrapper job 包装方法
func CronWrapper(j IJob) cron.Job {
	return CronRetryWrapper(j, 0)
}

// CronRetryWrapper job 重试包装方法
func CronRetryWrapper(j IJob, retry uint8) cron.Job {
	return &cronJob{
		job:      j,
		maxRetry: retry,
	}
}
