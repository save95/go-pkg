package job

import (
	"fmt"
	"log"
	"testing"
)

type customJob struct {
}

func newCustomJob() *customJob {
	return &customJob{}
}

func (c *customJob) Run() error {
	log.Print("run")

	return fmt.Errorf("some err")
}

func TestName(t *testing.T) {
	w := NewWrapper()
	w.Cron(newCustomJob()).Run()
}
