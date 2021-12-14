package timeutil

import (
	"time"

	util "github.com/save95/go-utils/timeutil"
)

// Format 格式化时间，失败返回空字符串
// 如果 date 为 nil，则返回空字符串
// Deprecated
func Format(date *time.Time, layout string) string {
	return util.Format(date, layout)
}
