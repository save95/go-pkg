package sliceutil

import "reflect"

func Is(arg interface{}) bool {
	val := reflect.ValueOf(arg)

	return val.Kind() == reflect.Slice
}
