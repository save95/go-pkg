package fileutil

import (
	"github.com/save95/go-utils/fsutil"
)

// Exist 判断文件是否存在
// Deprecated
func Exist(filename string) bool {
	return fsutil.Exist(filename)
}

// PathExist 判断路径/目录是否存在
// Deprecated
func PathExist(path string) bool {
	return fsutil.PathExist(path)
}

// Clear 清理目录下的所有文件
// Deprecated
func Clear(filePath string) error {
	return fsutil.Clear(filePath)
}

// Copy 复制文件
// Deprecated
func Copy(src, dst string) (int64, error) {
	return fsutil.Copy(src, dst)
}
