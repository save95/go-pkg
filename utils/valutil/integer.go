package valutil

import (
	util "github.com/save95/go-utils/valutil"
)

// Int 将任意值转成 int
// 如果输入值是 数字，则直接转换；
// 如果输入值是 boolean，则 true 转为 1；false 转为 0；
// 如果输入值是 string，则按字符串转换规则
// 否则抛出 ERROR
// Deprecated
func Int(any interface{}) (int, error) {
	return util.Int(any)
}
