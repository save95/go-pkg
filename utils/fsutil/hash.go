package fsutil

import (
	"io"

	newfsutil "github.com/save95/go-utils/fsutil"
)

// QiNiuHash 获得文件内容的 hash
// Deprecated
func QiNiuHash(f io.Reader, fsize int64) (etag string, err error) {
	return newfsutil.QiNiuHash(f, fsize)
}
