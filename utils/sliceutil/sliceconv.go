package sliceutil

import (
	"reflect"

	"github.com/save95/go-pkg/utils/valutil"
)

// ToAny 将任意切片转成 []interface{} 切片
func ToAny(slice interface{}) ([]interface{}, bool) {
	val := reflect.ValueOf(slice)

	if val.Kind() != reflect.Slice {
		return nil, false
	}

	sliceLen := val.Len()
	out := make([]interface{}, sliceLen)

	for i := 0; i < sliceLen; i++ {
		out[i] = val.Index(i).Interface()
	}

	return out, true
}

// ToFloat64 转成 []float64 切片，并过滤转换失败的数据
func ToFloat64(slice []interface{}) []float64 {
	if len(slice) == 0 {
		return []float64{}
	}

	res := make([]float64, 0, len(slice))
	for _, value := range slice {
		if v, ok := value.(float64); ok {
			res = append(res, v)
		}
	}

	return res
}

// ToInt 转成 []int 切片，并过滤转换失败的数据
func ToInt(slice []interface{}) []int {
	if len(slice) == 0 {
		return []int{}
	}

	res := make([]int, 0, len(slice))
	for _, value := range slice {
		if v, ok := value.(int); ok {
			res = append(res, v)
		}
	}

	return res
}

// ToPossibleInt 尽可能的转成 []int 切片
// 如：[]interface{}{1, "2", "b", "c", "", "3.0", "4.62"} => []int{1, 2, 3, 4}
func ToPossibleInt(slice []interface{}) []int {
	if len(slice) == 0 {
		return []int{}
	}

	res := make([]int, 0, len(slice))
	for _, value := range slice {
		if v, err := valutil.Int(value); err == nil {
			res = append(res, v)
		}
	}

	return res
}

// ToString 转成 []string 切片，并过滤转换失败的数据
func ToString(slice []interface{}) []string {
	if len(slice) == 0 {
		return []string{}
	}

	res := make([]string, 0, len(slice))
	for _, value := range slice {
		if v, ok := value.(string); ok {
			res = append(res, v)
		}
	}

	return res
}
