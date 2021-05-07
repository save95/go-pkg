package logger

import (
	"fmt"
	"os/exec"
	"sync"
	"testing"
	"time"

	"github.com/save95/xlog"
)

var (
	path = "log"
)

func TestLogger_NewLogger(t *testing.T) {
	logger := NewLogger(path, "", xlog.DailyStack)
	logger.Info("info log with args: ", "abc")
	logger.Info("info log with object: ", []string{"a", "b", "c"})
	logger.Infof("infof log %s ", "abc")
}

func TestLogger_rotate(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		logger := NewLogger(path, "rotate", xlog.DailyStack)
		for i := 0; i <= 80; i++ {
			logger.Info("test rotate logs")
			time.Sleep(time.Second)
		}
		wg.Done()
	}()

	go func() {
		for i := 1; i <= 3; i++ {
			time.Sleep(20 * time.Second)

			// 设置到系统
			now := time.Now()
			now.AddDate(0, 0, i)

			// 设置
			fmt.Println("set data...")
			cmd := exec.Command("date", "-s", now.Format("01/02/2006 15:04:05.999999999"))
			if err := cmd.Run(); nil != err {
				fmt.Printf("failed: %s\n\nPlease set time in system", err.Error())
			}
		}
		wg.Done()
	}()

	wg.Wait()
}
