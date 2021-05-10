package middleware

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/save95/go-pkg/http/types"
)

// Roles 角色权限中间件
func Roles(roles []types.IRole) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		svrCtx, err := types.ParserHttpContext(ctx)
		if nil != err {
			fmt.Println("role error: context convert failed")
			_ = ctx.AbortWithError(http.StatusForbidden, fmt.Errorf("context convert failed"))
			return
		}
		if !svrCtx.HasRole(roles) {
			fmt.Println("role error")
			_ = ctx.AbortWithError(http.StatusForbidden, fmt.Errorf("role error"))
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}
