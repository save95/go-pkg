package middleware

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// CORS 跨域处理
func CORS() gin.HandlerFunc {
	return cors.New(corsHandler{}.getCORSConfig())
}

type corsHandler struct {
}

func (ch corsHandler) getCORSConfig() cors.Config {
	return cors.Config{
		//AllowOrigins:     []string{"https://xxxx.com"},
		AllowOriginFunc: func(origin string) bool {
			//return origin == "https://xxxx.com"
			return true
		},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders: []string{
			"Origin", "Content-Type", "Accept", "User-Agent", "Cookie", "Authorization",
			"X-Auth-Token", "X-Token", "X-Requested-With",
			// https://www.npmjs.com/package/huge-uploader
			"uploader-chunk-number", "uploader-chunks-total", "uploader-file-id",
		},
		AllowCredentials: true,
		ExposeHeaders: []string{
			"Authorization", "Content-MD5",
			// 分页响应头
			"Link", "X-More-Resource", "X-Pagination-Info", "X-Total-Count",
		},
		MaxAge: 12 * time.Hour,
	}
}
