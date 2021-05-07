package jwt

import (
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/save95/go-pkg/constant"
	"github.com/save95/go-pkg/http/types"
)

// NewClaims 创建请求权
func NewClaims(user types.User) (*claims, error) {
	var jwtRoles []string
	for i := range user.Roles {
		jwtRoles = append(jwtRoles, user.Roles[i].String())
	}

	return &claims{
		Account: user.Name,
		UserID:  user.ID,
		Name:    user.Name,
		Roles:   jwtRoles,
	}, nil
}

// ParseClaims 从上下文中解析请求权
// 优先读取 http header 中的 X-Token 值；如果不存在，则读取 query string 中的 token 值
func ParseClaims(ctx *gin.Context) (*claims, error) {
	token := strings.TrimSpace(ctx.GetHeader(constant.HttpTokenHeaderKey))
	if len(token) == 0 {
		token, _ = ctx.GetQuery("token")
	}

	token = strings.TrimSpace(token)

	return parseToken(token)
}

func parseToken(token string) (*claims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if tokenClaims == nil {
		return nil, err
	}

	if claims, ok := tokenClaims.Claims.(*claims); ok && tokenClaims.Valid {
		return claims, nil
	}

	return nil, err
}
