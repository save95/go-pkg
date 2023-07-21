package middleware

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/save95/go-pkg/http/jwt"
	mjwt "github.com/save95/go-pkg/http/middleware/internal/jwt"
)

// JWTWith jwt 鉴权中间件
// 在用户登录成功后，配合 jwt.NewToken 生成 token
func JWTWith(opt *jwt.Option) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := mjwt.NewHandler(c, opt).Handle(); err != nil {
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
//
// usage:
//  ra := router.Group(
//		"/user",
//		middleware.JWTStatefulWith(
//			&jwt.Option{
//	    		RoleConvert:     NewRole,
//	    		RefreshDuration: 0, // 0-不自动刷新
//	    		Secret:          []byte(global.Config.App.Secret),
//	    	},
//			jwtstore.NewSingleRedisStore(global.SessionStoreClient), // 单地登录
//		),
//		middleware.Roles([]types.IRole{global.RoleBroker, global.RoleStar, global.RoleMember}),
//	)
func JWTStatefulWith(opt *jwt.Option, handler jwt.StatefulStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := mjwt.NewStatefulHandler(c, opt, handler).Handle(); err != nil {
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
func JWTStatefulWithout(opt *jwt.Option) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := mjwt.NewStatefulHandler(c, opt, nil).Handle(); err != nil {
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
