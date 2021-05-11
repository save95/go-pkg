package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/save95/go-pkg/constant"
	"github.com/save95/go-pkg/http/jwt"
	"github.com/save95/go-pkg/http/types"
)

// JWT jwt 鉴权中间件
func JWT(f types.ToRole) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := (userJwt{ctx: c, roleConvert: f}).handle(); err != nil {
			fmt.Println("Unauthorized")
			_ = c.AbortWithError(http.StatusUnauthorized, fmt.Errorf("unauthorized"))
			return
		}
		c.Next()
	}
}

type userJwt struct {
	ctx         *gin.Context
	roleConvert types.ToRole
}

// 鉴权处理
// 只负责验证是否登陆，不处理其他事务
func (uj userJwt) handle() error {
	claims, err := jwt.ParseClaims(uj.ctx)
	if err != nil {
		return errors.WithMessage(err, "登录信息错误, 请重新登录")
	}

	if time.Now().Unix() > claims.ExpiresAt {
		return errors.New("登录信息过期, 请重新登录")
	}

	uj.ctx.Set(constant.HttpTokenClaimsKey, claims)

	// 将 jwt 的用户角色转换
	roles, err := uj.rolesConvert(claims.Roles)
	if err != nil {
		return err
	}

	// 基础用户信息
	user := types.User{
		ID:    claims.UserID,
		Roles: roles,
		Name:  claims.Name,
	}

	// 写入自定义上下文
	if v, ok := uj.ctx.Get(constant.HttpCustomContextKey); ok {
		stx := v.(*types.HttpContext)
		stx.Set("user", user)
		stx.StorageTo(uj.ctx)
	}

	return nil
}

// 将 jwt 的用户角色转换
func (uj userJwt) rolesConvert(jwtRoles []string) ([]types.IRole, error) {
	var roles []types.IRole
	for i := range jwtRoles {
		r, err := uj.roleConvert(jwtRoles[i])
		if err != nil {
			return nil, err
		}
		roles = append(roles, r)
	}

	return roles, nil
}
