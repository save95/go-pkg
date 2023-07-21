package jwt

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/save95/go-pkg/http/types"
)

// MustParseJWTUser 解析 jwt User 信息。如果是静默模式（SilentMode=true）， User.ID 可能为零
func MustParseJWTUser(ctx *gin.Context, opt *Option) (*types.User, error) {
	if opt == nil || opt.RoleConvert == nil {
		return nil, errors.New("jwt option empty")
	}

	_, tk, err := ParseTokenWithSecret(ctx, opt.Secret)
	if nil != err {
		return nil, errors.WithMessage(err, "token error")
	}

	if tk.IsExpired() {
		return nil, errors.New("token expired")
	}

	// 基础用户信息
	user, err := tk.User(opt.RoleConvert)
	if err != nil {
		// 非静默模式，返回错误
		if !opt.SilentMode {
			return nil, err
		}
	}

	return user, nil
}
