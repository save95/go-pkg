package sliceutil

import (
	util "github.com/save95/go-utils/sliceutil"
)

// First 从切片中按下标取值
// 如果 any 不是切片，则返回 err
// 如果下标越界，则返回 err
// Deprecated
func First(slice interface{}, idx int) (interface{}, error) {
	return util.First(slice, idx)
}

// FirstString 从切片中按下标取值，并转化成字符串
// Deprecated
func FirstString(slice interface{}, idx int) (string, error) {
	return util.FirstString(slice, idx)
}
