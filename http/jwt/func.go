package jwt

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/save95/go-pkg/constant"
)

// ParseTokenWithGin 通过 gin.Context 初始化 token
// 从 gin.Context 优先读取 http header 中的 X-Token 值；如果不存在，则读取 query string 中的 token 值
func ParseTokenWithGin(ctx *gin.Context) (*token, error) {
	return ParseTokenWithGinSecret(ctx, jwtSecret)
}

// ParseTokenWithGinSecret 自定义 secret 初始化 token
func ParseTokenWithGinSecret(ctx *gin.Context, secret []byte) (*token, error) {
	c, err := parseClaims(ctx, secret)
	if nil != err {
		return nil, err
	}

	tk := newTokenWith(c).WithSecret(secret)

	return tk, nil
}

// parseClaims 从上下文中解析请求权
// 优先读取 http header 中的 X-Token 值；如果不存在，则读取 query string 中的 token 值
func parseClaims(ctx *gin.Context, secret []byte) (*claims, error) {
	tokenStr := strings.TrimSpace(ctx.GetHeader(constant.HttpTokenHeaderKey))
	if len(tokenStr) == 0 {
		tokenStr, _ = ctx.GetQuery("token")
	}

	tokenStr = strings.TrimSpace(tokenStr)

	c, err := parseToken(tokenStr, secret)
	if nil != err {
		return nil, err
	}

	c.IP = ctx.ClientIP()
	return c, nil
}

func parseToken(token string, secret []byte) (*claims, error) {
	if len(secret) == 0 {
		secret = jwtSecret
	}

	tokenClaims, err := jwt.ParseWithClaims(token, &claims{}, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return secret, nil
	})

	if tokenClaims == nil {
		return nil, err
	}

	if c, ok := tokenClaims.Claims.(*claims); ok && tokenClaims.Valid {
		return c, nil
	}

	return nil, err
}
