package fsutil

import (
	newfsutil "github.com/save95/go-utils/fsutil"
)

// Exist 判断文件是否存在
// Deprecated
func Exist(filename string) bool {
	return newfsutil.Exist(filename)
}

// IsDir 判断所给路径是否为文件夹
// Deprecated
func IsDir(path string) bool {
	return newfsutil.IsDir(path)
}

// PathExist 判断路径/目录是否存在
// Deprecated
func PathExist(path string) bool {
	return newfsutil.PathExist(path)
}

// Clear 清理目录下的所有文件
// Deprecated
func Clear(filePath string) error {
	return newfsutil.Clear(filePath)
}

// Copy 复制文件
// Deprecated
func Copy(src, dst string) (int64, error) {
	return newfsutil.CopyFile(src, dst, true)
}

// CopyFile 复制文件
// Deprecated
func CopyFile(src, dst string, overwrite bool) (int64, error) {
	return newfsutil.CopyFile(src, dst, overwrite)
}

// CopyPath 拷贝目录
// Deprecated
func CopyPath(src, dst string) error {
	return newfsutil.CopyPath(src, dst)
}
