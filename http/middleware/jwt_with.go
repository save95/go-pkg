package middleware

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/save95/go-pkg/http/jwt"
)

// JWTWith jwt 鉴权中间件
// 在用户登录成功后，配合 jwt.NewToken 生成 token
func JWTWith(opt *JWTOption) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := newJWTHandle(c, opt).handle(); err != nil {
			fmt.Printf("Unauthorized, slientMode=%v\n", opt.SilentMode)
			// 非静默模式，响应错误
			if !opt.SilentMode {
				_ = c.AbortWithError(http.StatusUnauthorized, fmt.Errorf("unauthorized"))
				return
			}
		}
		c.Next()
	}
}

// JWTStatefulWith 有状态的 jwt 鉴权中间件
// 需要配合 jwt.NewStatefulToken 使用（在用户登录成功后，调用该函数创建token）
func JWTStatefulWith(opt *JWTOption, handler jwt.StatefulStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := newJWTStatefulHandle(c, opt, handler).handle(); err != nil {
			fmt.Printf("Unauthorized, slientMode=%v\n", opt.SilentMode)
			// 非静默模式，响应错误
			if !opt.SilentMode {
				_ = c.AbortWithError(http.StatusUnauthorized, fmt.Errorf("unauthorized"))
				return
			}
		}
		c.Next()
	}
}

// JWTStatefulWithout 有状态的 jwt 鉴权中间件，仅校验 jwt 是否合法，不校验状态
// 需要配合 jwt.NewStatefulToken 使用（在用户登录成功后，调用该函数创建token）
func JWTStatefulWithout(opt *JWTOption) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := newJWTStatefulHandle(c, opt, nil).handle(); err != nil {
			fmt.Printf("Unauthorized, slientMode=%v\n", opt.SilentMode)
			// 非静默模式，响应错误
			if !opt.SilentMode {
				_ = c.AbortWithError(http.StatusUnauthorized, fmt.Errorf("unauthorized"))
				return
			}
		}
		c.Next()
	}
}
