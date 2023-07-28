package job

import "github.com/robfig/cron/v3"

// IJob job 约定
// Deprecated
type IJob interface {
	Run() error
}

type ICommandJob interface {
	Run(args ...string) error
}

type IWrapper interface {
	FromCommandJob(job ICommandJob, args ...string) cron.Job
}
