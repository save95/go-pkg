package logger

import (
	"os"
	"path"
	"runtime"
)

func getDefaultDir() string {
	if runtime.GOOS == "linux" {
		return defaultDir
	}

	return path.Join(os.TempDir(), "logs", "go-pkg")
}
