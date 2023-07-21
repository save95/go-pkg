package types

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"math/rand"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/save95/go-pkg/constant"
	"github.com/save95/go-pkg/framework/logger"
)

// NewHttpContext 创建自定义上下文
func NewHttpContext() *HttpContext {
	traceId := traceId()
	return &HttpContext{
		traceId: traceId,
		logger:  logger.NewDefaultTraceLogger(traceId),
	}
}

// ParserHttpContext 从 gin 上下文中解析自定义上下文
// Deprecated. use MustParseHttpContext
func ParserHttpContext(ctx context.Context) (*HttpContext, error) {
	return MustParseHttpContext(ctx)
}

// MustParseHttpContext 从 gin 上下文中解析自定义上下文
func MustParseHttpContext(ctx context.Context) (*HttpContext, error) {
	gtx, ok := ctx.(*gin.Context)
	if !ok {
		return nil, errors.New("to GinContext failed")
	}

	v, ok := gtx.Get(constant.HttpCustomContextKey)
	if !ok {
		return nil, errors.New("get HttpCustomContext failed")
	}

	rtx, ok := v.(*HttpContext)
	if !ok {
		return nil, errors.New("to HttpContext failed")
	}

	return rtx, nil
}

// NOTE 也可以使用三方UUID
func traceId() string {
	h := md5.New()
	rand.Seed(time.Now().UnixNano())
	h.Write([]byte(strconv.FormatInt(rand.Int63(), 10)))
	h.Write([]byte("-"))
	h.Write([]byte(strconv.FormatInt(time.Now().UnixNano(), 10)))
	h.Write([]byte("-"))
	h.Write([]byte(strconv.FormatInt(int64(rand.Int31()), 10)))
	return hex.EncodeToString(h.Sum([]byte("server-api")))
}
