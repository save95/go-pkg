package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/save95/go-pkg/http/middleware/internal/logger"
	"github.com/save95/xlog"
)

var _otherLogHandlers = []string{
	"github.com/save95/go-pkg/http/middleware.HttpPrinter",
}

type HttpLoggerOption struct {
	Logger    xlog.XLog
	OnlyError bool // 仅发生错误时，打印日志；否则，打印所有请求
}

// HttpLogger http 日志中间件；
// 如果有其他内置日志，则该中间件不操作；内置日志有: HttpPrinter 等
//
// usage:
//   r.Use(middleware.HttpLogger(middleware.HttpLoggerOption{
//	 	Logger:    global.Log,
//	 	OnlyError: global.Config.Log.HttpLogOnlyError,
//	 }))
func HttpLogger(opt HttpLoggerOption) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 注入 gin.ResponseWriter
		l := logger.New(c)

		c.Next()

		// 是否只打印错误日志
		needLog := true
		errors := c.Errors.ByType(gin.ErrorTypeAny)
		if len(errors) == 0 && opt.OnlyError {
			needLog = false
		}

		// 是否有其他内置日志中间件，有则不打印
		hasOtherLogger := false
		for _, s := range _otherLogHandlers {
			if strings.Contains(strings.Join(c.HandlerNames(), ", "), s) {
				hasOtherLogger = true
				break
			}
		}

		if needLog && !hasOtherLogger {
			opt.Logger.Info(l.String())
		}
	}
}
