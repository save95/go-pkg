package fsutil

import (
	newfsutil "github.com/save95/go-utils/fsutil"
)

// BlockRead 分块读取大文件，每次读取 4M 内容
// Deprecated
func BlockRead(fileName string, handle func([]byte) error) error {
	return newfsutil.BlockRead(fileName, handle)
}
