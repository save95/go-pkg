package middleware

import (
	"bytes"
	"io/ioutil"

	"github.com/gin-gonic/gin"
	"github.com/save95/go-pkg/constant"
	"github.com/save95/go-pkg/http/types"
)

// HttpContext 注入自定义上下文
func HttpContext() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		stx := types.NewHttpContext()

		// 取出 request body
		body, _ := ctx.GetRawData()
		stx.Set(constant.HttpCustomRawRequestBodyKey, body)

		// 重新写入 request body
		ctx.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body))

		// 注册自定义上下文
		ctx.Set(constant.HttpCustomContextKey, stx)

		ctx.Next()
	}
}
