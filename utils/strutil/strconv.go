package strutil

import (
	"strconv"
	"time"
)

// ToInt 字符串转成 int，失败返回0
func ToInt(str string) int {
	i, err := strconv.Atoi(str)
	if nil != err {
		return 0
	}

	return i
}

// ToIntWith 字符串转成 int，失败返回默认值
func ToIntWith(str string, defaultValue int) int {
	i := ToInt(str)
	if i == 0 {
		return defaultValue
	}

	return i
}

// ToTime 字符串转 time，失败返回 nil
func ToTime(str string) *time.Time {
	return ToTimeWithLayout(str, "2006-01-02 15:04:05")
}

// ToTimeWithLayout 按格式化转成成 time，失败返回 nil
func ToTimeWithLayout(str string, layout string) *time.Time {
	t, err := time.ParseInLocation(layout, str, time.Local)
	if nil != err {
		return nil
	}

	return &t
}
