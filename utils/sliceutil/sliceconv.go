package sliceutil

import (
	util "github.com/save95/go-utils/sliceutil"
)

// ToAny 将任意切片转成 []interface{} 切片
// Deprecated
func ToAny(slice interface{}) ([]interface{}, bool) {
	return util.ToAny(slice)
}

// ToFloat64 转成 []float64 切片，并过滤转换失败的数据
// Deprecated
func ToFloat64(slice []interface{}) []float64 {
	return util.ToFloat64(slice)
}

// ToInt 转成 []int 切片，并过滤转换失败的数据
// Deprecated
func ToInt(slice []interface{}) []int {
	return util.ToInt(slice)
}

// ToPossibleInt 尽可能的转成 []int 切片
// 如：[]interface{}{1, "2", "b", "c", "", "3.0", "4.62"} => []int{1, 2, 3, 4}
// Deprecated
func ToPossibleInt(slice []interface{}) []int {
	return util.ToPossibleInt(slice)
}

// ToString 转成 []string 切片，并过滤转换失败的数据
// Deprecated
func ToString(slice []interface{}) []string {
	return util.ToString(slice)
}
