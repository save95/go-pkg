package fsutil

import (
	newfsutil "github.com/save95/go-utils/fsutil"
)

// Download 下载远程文件
// Deprecated
func Download(filename, url string) error {
	return newfsutil.Download(filename, url)
}
