package middleware

import (
	"github.com/gin-gonic/gin"
	r "github.com/save95/go-pkg/http/restful"
	"github.com/save95/xlog"
)

// HttpPrinter 打印 http 信息中间件；展示 request / response 等信息
func HttpPrinter(logger xlog.XLog) gin.HandlerFunc {
	return func(c *gin.Context) {
		l := r.NewHttpLogger(c)

		c.Next()

		logger.Info(l.String())
	}
}
