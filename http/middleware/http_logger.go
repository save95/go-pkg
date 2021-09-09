package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	r "github.com/save95/go-pkg/http/restful"
)

var _otherLogHandlers = []string{
	"github.com/save95/go-pkg/http/middleware.HttpPrinter",
}

// HttpLogger http 日志中间件；
// 如果有其他内置日志，则该中间件不操作；内置日志有: HttpPrinter 等
func HttpLogger(opt r.LogOption) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 注入 gin.ResponseWriter
		l := r.NewHttpLogger(c)

		c.Next()

		// 是否只打印错误日志
		needLog := true
		errors := c.Errors.ByType(gin.ErrorTypeAny)
		if len(errors) == 0 && opt.OnlyError {
			needLog = false
		}

		// 是否有其他内置日志中间件，有则不打印
		hasOtherLog := false
		for _, s := range _otherLogHandlers {
			if strings.Contains(strings.Join(c.HandlerNames(), ", "), s) {
				hasOtherLog = true
				break
			}
		}

		if needLog && !hasOtherLog {
			opt.Logger.Info(l.String())
		}
	}
}
