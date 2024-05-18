package job

import "github.com/robfig/cron/v3"

type ICommandJob interface {
	Run(args ...string) error
}

type IWrapper interface {
	FromCommandJob(job ICommandJob, args ...string) cron.Job
}

type ICronjobRegister interface {
	Register(spec string, cmd ICommandJob)
}

type ICommandRegister interface {
	Register(name string, cmd ICommandJob)
}
