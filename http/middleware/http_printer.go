package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/save95/go-pkg/http/middleware/internal/logger"
	"github.com/save95/xlog"
)

// HttpPrinter 打印 http 信息中间件；展示 request / response 等信息
//
// usage:
//   r.Use(middleware.HttpContext(global.Log))
//
//   router.Any("/endpoint", middleware.HttpPrinter(global.Log), ping.Controller{}.Endpoint)
func HttpPrinter(log xlog.XLog) gin.HandlerFunc {
	return func(c *gin.Context) {
		l := logger.New(c)

		c.Next()

		log.Info(l.String())
	}
}
