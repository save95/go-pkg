package fsutil

import (
	"fmt"
	"io"
	"os"
)

// Exist 判断文件是否存在
func Exist(filename string) bool {
	// os.Stat获取文件信息
	if _, err := os.Stat(filename); err != nil {
		return false
	}
	return true
}

// PathExist 判断路径/目录是否存在
func PathExist(path string) bool {
	// os.Stat获取文件信息
	if _, err := os.Stat(path); err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

// Clear 清理目录下的所有文件
func Clear(filePath string) error {
	if Exist(filePath) {
		if err := os.RemoveAll(filePath); nil != err {
			return err
		}
	}

	return nil
}

// Copy 复制文件
func Copy(src, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("`%s` is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer func() {
		_ = source.Close()
	}()

	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer func() {
		_ = destination.Close()
	}()

	return io.Copy(destination, source)
}
