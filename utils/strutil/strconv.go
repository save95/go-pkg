package strutil

import (
	"time"

	util "github.com/save95/go-utils/strutil"
)

// ToInt 字符串转成 int，失败返回0
// Deprecated
func ToInt(str string) int {
	return util.ToInt(str)
}

// ToIntWith 字符串转成 int，失败返回默认值
// Deprecated
func ToIntWith(str string, defaultValue int) int {
	return util.ToIntWith(str, defaultValue)
}

// ToTime 字符串转 time，失败返回 nil
// Deprecated
func ToTime(str string) *time.Time {
	return util.ToTimeWithLayout(str, "2006-01-02 15:04:05")
}

// ToTimeWithLayout 按格式化转成成 time，失败返回 nil
// Deprecated
func ToTimeWithLayout(str string, layout string) *time.Time {
	return util.ToTimeWithLayout(str, layout)
}
