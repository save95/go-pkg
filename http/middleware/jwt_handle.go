package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/save95/go-pkg/constant"
	"github.com/save95/go-pkg/http/jwt"
	"github.com/save95/go-pkg/http/types"
)

type jwtHandle struct {
	ctx *gin.Context
	opt *JWTOption
}

func newJWTHandle(ctx *gin.Context, opt *JWTOption) *jwtHandle {
	return &jwtHandle{
		ctx: ctx,
		opt: opt,
	}
}

// 鉴权处理
// 只负责验证是否登陆，不处理其他事务
func (h *jwtHandle) handle() error {
	if h.opt == nil || h.opt.RoleConvert == nil {
		return errors.New("jwt option empty")
	}

	token, err := jwt.ParseTokenWithGinSecret(h.ctx, h.opt.Secret)
	if nil != err {
		return errors.WithMessage(err, "token error")
	}

	if token.IsExpired() {
		return errors.New("token expired")
	}

	if token.IsStateful() {
		return errors.New("token is stateful, please use middleware.JWTStatefulWith")
	}

	// 自动刷新 token
	if h.opt.RefreshDuration > 0 {
		token.RefreshNear(h.opt.RefreshDuration)
		// 失败，则跳过，只处理成功的情况
		if newToken, err := token.ToString(); nil == err {
			h.ctx.Header(constant.HttpTokenHeaderKey, newToken)
		}
	}

	// 基础用户信息
	user, err := token.User(h.opt.RoleConvert)
	if err != nil {
		if h.opt.SilentMode {
			return nil
		} else {
			return err
		}
	}

	// 写入自定义上下文
	if v, ok := h.ctx.Get(constant.HttpCustomContextKey); ok {
		stx := v.(*types.HttpContext)
		stx.Set("user", *user)
		stx.StorageTo(h.ctx)
	}

	return nil
}
