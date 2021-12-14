package valutil

import (
	util "github.com/save95/go-utils/valutil"
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
// Deprecated
func Bool(any interface{}) (bool, error) {
	return util.Bool(any)
}
