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
		svrCtx, err := types.MustParseHttpContext(ctx)
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

// WithRole 角色权限中间件
func WithRole(role types.IRole, roles ...types.IRole) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		svrCtx, err := types.MustParseHttpContext(ctx)
		if nil != err {
			fmt.Println("role error: context convert failed")
			_ = ctx.AbortWithError(http.StatusForbidden, fmt.Errorf("context convert failed"))
			return
		}

		rs := []types.IRole{role}
		if len(roles) > 0 {
			rs = append(rs, roles...)
		}
		if !svrCtx.HasRole(rs) {
			fmt.Println("role error")
			_ = ctx.AbortWithError(http.StatusForbidden, fmt.Errorf("role error"))
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}

// RoleFunc 角色控制器中间件。
// 如果用户满足指定角色要求，则使用调用 action，并在完成后进入下一个中间件；
// 如果用户不满足指定角色要求，则直接进入下一个中间件
func RoleFunc(action gin.HandlerFunc, roles ...types.IRole) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if htx, err := types.MustParseHttpContext(ctx); nil == err && htx.HasRole(roles) {
			action(ctx)
		}

		ctx.Next()
	}
}

// RoleFuncAbort 角色控制器独占中间件。
// 如果用户符合指定角色，则使用调用 action，并在完成后进入下一个中间件；
// 如果用户不满足指定角色要求，则中断链路，返回 http status 400 错误
func RoleFuncAbort(action gin.HandlerFunc, roles ...types.IRole) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if htx, err := types.MustParseHttpContext(ctx); nil == err && htx.HasRole(roles) {
			action(ctx)

			ctx.Next()
		} else {
			fmt.Println("role error, abort")
			ctx.AbortWithStatus(http.StatusForbidden)
		}
	}
}
