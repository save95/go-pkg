package restful

import (
	"errors"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/save95/go-pkg/http/types"
)

type handler struct {
	ctx     *gin.Context
	version types.ApiVersion
}

func New(ctx *gin.Context, version types.ApiVersion) *handler {
	return &handler{
		ctx:     ctx,
		version: version,
	}
}

func (h handler) Handle() error {
	if err := h.parseAccept(); nil != err {
		return err
	}

	return nil
}

func (h handler) parseAccept() error {
	stx, err := types.MustParseHttpContext(h.ctx)
	if nil != err {
		return err
	}

	// see: https://developer.github.com/v3/media/#request-specific-version
	// application/vnd.server[.version].param[+json]
	// eg: application/vnd.server.v1.raw+json
	accept := h.ctx.GetHeader("Accept")

	// 默认值
	if len(accept) == 0 || accept == "*/*" || strings.Contains(accept, "application/json") {
		stx.Set("version", h.version)
		stx.Set("bodyProperty", types.BodyPropertyRaw)
		stx.StorageTo(h.ctx)
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
		stx.StorageTo(h.ctx)
		return nil
	}

	return errors.New("not support custom media type")
}
