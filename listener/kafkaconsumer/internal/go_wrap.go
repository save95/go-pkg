package internal

import (
	"runtime/debug"
	"time"
)

// goWrap 开启一个 goroutine。保证采用一致的 recover 和 重启机制
func goWrap(handler *groupHandlerParam, fn func(handler *groupHandlerParam)) {
	go func() {
		// FIX: panic(nil)
		paniked := true

		defer func() {
			if v := recover(); paniked && v != nil {
				debug.PrintStack()
				// 延迟5s再重启
				time.Sleep(5 * time.Second)
				goWrap(handler, fn)
			}
		}()

		fn(handler)

		paniked = false
	}()
}
