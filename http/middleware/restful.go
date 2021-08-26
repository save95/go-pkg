package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"
	"github.com/save95/go-pkg/http/types"
)

// RESTFul Restful 标准检测解析中间件
func RESTFul(version types.ApiVersion) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if err := (restful{ctx: ctx, version: version}).Handle(); nil != err {
			fmt.Printf("not support accept: %s\n", err)
			_ = ctx.AbortWithError(http.StatusBadRequest, fmt.Errorf("not support accept"))
			return
		}

		ctx.Next()
	}
}

type restful struct {
	ctx     *gin.Context
	version types.ApiVersion
}

func (r restful) Handle() error {
	if err := r.parseAccept(); nil != err {
		return err
	}

	return nil
}

func (r restful) parseAccept() error {
	// see: https://developer.github.com/v3/media/#request-specific-version
	// application/vnd.server[.version].param[+json]
	// eg: application/vnd.server.v1.raw+json
	accept := r.ctx.GetHeader("Accept")
	if len(accept) == 0 {
		return errors.New("header error")
	}

	stx, err := types.ParserHttpContext(r.ctx)
	if nil != err {
		return err
	}

	// 默认值
	if accept == "*/*" || accept == "application/json" {
		stx.Set("version", r.version)
		stx.Set("bodyProperty", types.BodyPropertyRaw)
		stx.StorageTo(r.ctx)
		return nil
	}

	// 解析自定义媒体类型
	re := regexp.MustCompile(`application/vnd\.server(\.(v\S+?))(\.(raw|text|html|full))?\+json`)
	params := re.FindStringSubmatch(accept)
	//fmt.Printf("accept: %+v\n  %+v\n", accept, params)
	if len(params) == 5 {
		av := types.ApiVersion(params[2])
		bp := types.BodyProperty(params[4])
		if !av.Verify() || !bp.Verify() {
			return errors.New("not support custom media type")
		}

		stx.Set("version", av)
		stx.Set("bodyProperty", bp)
		stx.StorageTo(r.ctx)
		return nil
	}

	return errors.New("not support custom media type")
}
