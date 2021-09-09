package restful

import "github.com/save95/xlog"

type LogOption struct {
	Logger    xlog.XLog
	OnlyError bool // 仅发生错误时，打印日志；否则，打印所有请求
}
