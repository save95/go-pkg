package jwt

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

var jwtSecret = []byte("go-pkg.JwtSecret")

type claims struct {
	jwt.StandardClaims

	Account string   `json:"account,omitempty"` // 账号
	UserID  uint     `json:"uid"`               // 用户ID
	Name    string   `json:"name"`              // 姓名
	Roles   []string `json:"roles"`             // 角色组
}

// RefreshToken 刷新 token
func (c *claims) RefreshToken() {
	nowTime := time.Now()
	expireTime := nowTime.Add(24 * time.Hour)
	c.StandardClaims.ExpiresAt = expireTime.Unix()
}

// ToTokenString 转成 token 字符串
func (c *claims) ToTokenString() (string, error) {
	c.RefreshToken()
	c.StandardClaims.Issuer = "go-pkg"
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, c)

	return tokenClaims.SignedString(jwtSecret)
}
