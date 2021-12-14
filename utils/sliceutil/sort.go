package sliceutil

import (
	util "github.com/save95/go-utils/sliceutil"
)

// Shuffle 随机打散切片
// Deprecated
func Shuffle(slice []interface{}) {
	util.Shuffle(slice)
}
