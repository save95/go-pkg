package middleware

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// JWTWith jwt 鉴权中间件
func JWTWith(opt *JWTOption) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := newJWTHandle(c, opt).handle(); err != nil {
			fmt.Println("Unauthorized")
			_ = c.AbortWithError(http.StatusUnauthorized, fmt.Errorf("unauthorized"))
			return
		}
		c.Next()
	}
}
