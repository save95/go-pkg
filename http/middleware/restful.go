package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/save95/go-pkg/http/middleware/internal/restful"
	"github.com/save95/go-pkg/http/types"
)

// RESTFul Restful 标准检测解析中间件
func RESTFul(version types.ApiVersion) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if err := restful.New(ctx, version).Handle(); nil != err {
			fmt.Printf("not support accept: %s\n", err)
			_ = ctx.AbortWithError(http.StatusBadRequest, fmt.Errorf("not support accept"))
			return
		}

		ctx.Next()
	}
}

type IgnorePath struct {
	Path   string
	Method string
}

// RESTFulWithIgnores 忽略指定 path 的Restful 标准检测解析中间件
// 一般，用在部分直接下载或浏览器直接访问的接口
func RESTFulWithIgnores(version types.ApiVersion, ignorePaths ...IgnorePath) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		for _, ignore := range ignorePaths {
			if ignore.Path == ctx.FullPath() &&
				strings.ToLower(ctx.Request.Method) == strings.ToLower(ignore.Method) {
				ctx.Next()
				return
			}
		}

		if err := restful.New(ctx, version).Handle(); nil != err {
			fmt.Printf("not support accept: %s\n", err)
			_ = ctx.AbortWithError(http.StatusBadRequest, fmt.Errorf("not support accept"))
			return
		}

		ctx.Next()
	}
}
