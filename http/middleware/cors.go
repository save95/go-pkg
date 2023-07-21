package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/save95/go-pkg/http/middleware/internal/cors"
)

// CORS 跨域处理
//
// usage:
//   r.Use(middleware.CORS())
//
// 	 r.Use(middleware.CORS(
//	 	middleware.WithCORSAllowOriginFunc(func(origin string) bool {
//	 		//return origin == "https://xxxx.com"
//	 		return true
//	 	}),
//	 	middleware.WithCORSAllowHeaders("X-Custom-Key"),
//	 	middleware.WithCORSExposeHeaders("X-Custom-Key"),
//	 	middleware.WithCORSMaxAge(24*time.Hour),
//	 ))
func CORS(opts ...cors.Option) gin.HandlerFunc {
	return cors.New(opts...)
}

func WithCORSAllowOriginFunc(fun func(origin string) bool) cors.Option {
	return cors.WithAllowOriginFunc(fun)
}

func WithCORSAllowMethods(methods ...string) cors.Option {
	return cors.WithAllowMethods(methods...)
}

func WithCORSAllowHeaders(keys ...string) cors.Option {
	return cors.WithAllowHeaders(keys...)
}

func WithCORSExposeHeaders(keys ...string) cors.Option {
	return cors.WithExposeHeaders(keys...)
}

func WithCORSMaxAge(d time.Duration) cors.Option {
	return cors.WithMaxAge(d)
}
