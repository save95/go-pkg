package sliceutil

import util "github.com/save95/go-utils/sliceutil"

// Get 从切片中按下标取值
// 如果 any 不是切片，则返回 nil
// 如果下标越界，则返回 nil
// 如果需要返回错误，请使用 First
// Deprecated
func Get(slice interface{}, idx int) interface{} {
	return util.Get(slice, idx)
}

// GetString 从切片中按下标取值，并转化成字符串
// 如果遇到错误，则返回空字符串，如果需要关注错误，请使用 FirstString
// Deprecated
func GetString(slice interface{}, idx int) string {
	return util.GetString(slice, idx)
}
