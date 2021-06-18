package sliceutil

import (
	"fmt"
)

// First 从切片中按下标取值
// 如果 any 不是切片，则返回 err
// 如果下标越界，则返回 err
func First(slice interface{}, idx int) (interface{}, error) {
	if !Is(slice) {
		return nil, fmt.Errorf("input not slice")
	}

	cols, ok := ToAny(slice)
	if !ok {
		return nil, fmt.Errorf("slice item convert to interface failed")
	}

	if len(cols) <= idx || idx < 0 {
		return nil, fmt.Errorf("`idx` out of bounds")
	}

	return cols[idx], nil
}

// FirstString 从切片中按下标取值，并转化成字符串
func FirstString(slice interface{}, idx int) (string, error) {
	val, err := First(slice, idx)
	if nil != err {
		return "", err
	}

	str, ok := val.(string)
	if !ok {
		return "", fmt.Errorf("slice item convert to string failed")
	}

	return str, nil
}
