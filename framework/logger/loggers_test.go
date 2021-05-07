package logger

import (
	"strconv"
	"sync"
	"testing"

	"github.com/save95/xlog"
)

var (
	basePath   = "log"
	categories = []string{
		"service",
		"payment",
		"order",
	}
	userCategory = "user"
)

func TestNewLoggers(t *testing.T) {
	ls := NewLoggers(basePath, categories, xlog.DailyStack)
	ls.GetLogger("user1").Info("user1 log")

	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		for i := 1; i < 10; i++ {
			ls.GetLogger(userCategory).Info("user log " + strconv.Itoa(i))
		}
		wg.Done()
	}()

	go func() {
		for i := 1; i < 10; i++ {
			ls.GetLogger(userCategory).Info("user log2 " + strconv.Itoa(i))
		}
		wg.Done()
	}()

	wg.Wait()
}
