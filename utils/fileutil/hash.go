package fileutil

import (
	"io"

	"github.com/save95/go-utils/fsutil"
)

// QiNiuHash 获得文件内容的 hash
// see: fsutil.QiNiuHash()
// Deprecated
func QiNiuHash(f io.Reader, fsize int64) (etag string, err error) {
	return fsutil.QiNiuHash(f, fsize)
}
