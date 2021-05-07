package valutil

import (
	"errors"
	"strings"
)

// Bool 将任意值转成 bool
// 如果传入值是 boolean，直接强制转换返回；
// 如果传入值是 string，则按以下规则转换：
//    "true, yes" -> true
//    "false, no" -> false
//    "" -> false
//    "0.0···001 ... 1 ... ∞" -> true
//    "-∞ ... -1 ... -0.1 ... 0" -> false
//    "other word" -> ERROR
// 如果传入值是 数字，大于零返回 true，否则返回 false
//    0.0···001 ... 1 ... ∞ -> true
//    -∞ ... -1 ... -0.1 ... 0 -> false
// 如果传入值是其他类型，则返回 ERROR
func Bool(any interface{}) (bool, error) {
	if v, ok := any.(bool); ok {
		return v, nil
	}

	if v, ok := any.(string); ok {
		switch strings.ToLower(v) {
		case "true", "yes":
			return true, nil
		case "false", "no", "":
			return false, nil
		default:
			i, err := Int(v)
			if nil != err {
				return false, errors.New("is string, but cannot convert it")
			}

			if i > 0 {
				return true, nil
			} else {
				return false, nil
			}
		}
	}

	i, err := Int(any)
	if nil == err {
		if i > 0 {
			return true, nil
		} else {
			return false, nil
		}
	}

	return false, errors.New("cannot convert it")
}
