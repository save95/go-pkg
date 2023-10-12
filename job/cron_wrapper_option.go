package job

import (
	"context"

	"github.com/save95/xlog"
)

type WrapperOption func(*cronJobWrapper)

func WrapWithContext(ctx context.Context) WrapperOption {
	return func(job *cronJobWrapper) {
		job.ctx = ctx
	}
}

func WrapWithLogger(log xlog.XLogger) WrapperOption {
	return func(job *cronJobWrapper) {
		job.log = log
	}
}

func WrapWithMaxRetry(retry uint8) WrapperOption {
	return func(job *cronJobWrapper) {
		job.maxRetry = retry
	}
}

func WrapWithFailedSaver(saver func(jobName string, in []string, err error)) WrapperOption {
	return func(job *cronJobWrapper) {
		job.failedSaver = saver
	}
}
