package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/save95/go-pkg/constant"
	"github.com/save95/go-pkg/http/jwt"
	"github.com/save95/go-pkg/http/types"
)

type jwtStatefulHandle struct {
	ctx *gin.Context
	opt *JWTOption

	store jwt.StatefulStore // token 状态处理器
}

func newJWTStatefulHandle(ctx *gin.Context, opt *JWTOption, store jwt.StatefulStore) *jwtStatefulHandle {
	return &jwtStatefulHandle{
		ctx:   ctx,
		opt:   opt,
		store: store,
	}
}

// 鉴权处理
// 只负责验证是否登陆，不处理其他事务
func (h *jwtStatefulHandle) handle() error {
	if h.opt == nil || h.opt.RoleConvert == nil {
		return errors.New("jwt option empty")
	}

	if h.store == nil || h.store.Check == nil {
		return errors.New("token is stateful, but checker undefined")
	}

	tokenStr, token, err := jwt.ParseStatefulTokenWithGinSecret(h.ctx, h.opt.Secret)
	if nil != err {
		return errors.WithMessage(err, "token error")
	}

	if token.IsExpired() {
		return errors.New("token expired")
	}

	if !token.IsStateful() {
		return errors.New("token is not stateful, please use middleware.JWTWith")
	}

	// 基础用户信息
	user, err := token.User(h.opt.RoleConvert)
	if err != nil {
		return err
	}

	// 判断 jwt 是否为有状态，通过函数处理判断状态是否有效
	if err := h.store.Check(user.GetID(), tokenStr); nil != err {
		return err
	}

	// 自动刷新 token
	if h.opt.RefreshDuration > 0 {
		token.RefreshNear(h.opt.RefreshDuration)
		// 失败，则跳过，只处理成功的情况
		if newToken, err := token.ToString(); nil == err {
			h.ctx.Header(constant.HttpTokenHeaderKey, newToken)
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
