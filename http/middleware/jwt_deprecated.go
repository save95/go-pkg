package middleware

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/save95/go-pkg/http/types"
)

// JWT jwt 鉴权中间件
// Deprecated
// 请使用 JWTWith 替代
func JWT(f types.ToRole) gin.HandlerFunc {
	return func(c *gin.Context) {
		opt := &JWTOption{
			RoleConvert:     f,
			RefreshDuration: 0,
		}
		if err := newJWTHandle(c, opt).handle(); err != nil {
			fmt.Println("Unauthorized")
			_ = c.AbortWithError(http.StatusUnauthorized, fmt.Errorf("unauthorized"))
			return
		}
		c.Next()
	}
}
