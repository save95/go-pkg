package ziputil

import (
	util "github.com/save95/go-utils/ziputil"
)

// CompressPath 压缩一个指定目录
// Deprecated
func CompressPath(dst, src string) error {
	return util.CompressPath(dst, src)
}
