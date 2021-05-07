package valutil

import (
	"errors"
	"strconv"
)

// Int 将任意值转成 int
// 如果输入值是 数字，则直接转换；
// 如果输入值是 boolean，则 true 转为 1；false 转为 0；
// 如果输入值是 string，则按字符串转换规则
// 否则抛出 ERROR
func Int(any interface{}) (int, error) {
	if v, ok := any.(int); ok {
		return v, nil
	}
	if v, ok := any.(uint); ok {
		return int(v), nil
	}

	if v, ok := any.(int64); ok {
		return int(v), nil
	}
	if v, ok := any.(uint64); ok {
		return int(v), nil
	}

	if v, ok := any.(int32); ok {
		return int(v), nil
	}
	if v, ok := any.(uint32); ok {
		return int(v), nil
	}

	if v, ok := any.(int16); ok {
		return int(v), nil
	}
	if v, ok := any.(uint16); ok {
		return int(v), nil
	}

	if v, ok := any.(int8); ok {
		return int(v), nil
	}
	if v, ok := any.(uint8); ok {
		return int(v), nil
	}

	if v, ok := any.(float64); ok {
		return int(v), nil
	}
	if v, ok := any.(float32); ok {
		return int(v), nil
	}

	if v, ok := any.(string); ok {
		if vv, err := strconv.Atoi(v); err == nil {
			return vv, nil
		}

		if fv, err := strconv.ParseFloat(v, 64); err == nil {
			return int(fv), nil
		}
	}

	if v, ok := any.(bool); ok {
		if v {
			return 1, nil
		} else {
			return 0, nil
		}
	}

	return 0, errors.New("interface convert to int failed")
}
