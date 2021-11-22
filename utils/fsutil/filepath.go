package fsutil

import (
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// Exist 判断文件是否存在
func Exist(filename string) bool {
	info, err := os.Stat(filename)
	return !os.IsNotExist(err) && !info.IsDir()
}

// IsDir 判断所给路径是否为文件夹
func IsDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}

// PathExist 判断路径/目录是否存在
func PathExist(path string) bool {
	info, _ := os.Stat(path)
	return info != nil && info.IsDir()
}

// Clear 清理目录下的所有文件
func Clear(filePath string) error {
	if Exist(filePath) || PathExist(filePath) {
		if err := os.RemoveAll(filePath); nil != err {
			return err
		}
	}

	return nil
}

// Copy 复制文件
func Copy(src, dst string) (int64, error) {
	return CopyFile(src, dst, true)
}

// CopyFile 复制文件
func CopyFile(src, dst string, overwrite bool) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() || sourceFileStat.IsDir() {
		return 0, fmt.Errorf("`%s` is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer func() {
		_ = source.Close()
	}()

	// 如果目标目录不存在，则创建
	dstPath := path.Dir(dst)
	if !PathExist(dstPath) {
		sourcePathStat, _ := os.Stat(path.Dir(src))
		if sourcePathStat != nil && !sourcePathStat.IsDir() {
			return 0, fmt.Errorf("`%s` file dir stat error", src)
		}

		_ = os.MkdirAll(dstPath, sourcePathStat.Mode())
	}

	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer func() {
		_ = destination.Close()
	}()

	return io.Copy(destination, source)
}

// CopyPath 拷贝目录
func CopyPath(src, dst string) error {
	if info, err := os.Stat(src); os.IsNotExist(err) || !info.IsDir() {
		return fmt.Errorf("`src` not dir")
	}

	return filepath.Walk(src, func(filename string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relationPath := strings.Trim(strings.Replace(filename, src, "", 1), string(os.PathSeparator))
		dstFilename := path.Join(dst, relationPath)
		if !info.IsDir() {
			if _, err := Copy(filename, dstFilename); nil != err {
				return err
			}
		}
		return nil
	})
}
