package timeutil

import "time"

// Format 格式化时间，失败返回空字符串
// 如果 date 为 nil，则返回空字符串
func Format(date *time.Time, layout string) string {
	if date != nil {
		return date.Format(layout)
	}

	return ""
}
