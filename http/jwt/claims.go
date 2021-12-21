package jwt

import (
	"github.com/golang-jwt/jwt/v4"
)

var jwtSecret = []byte("go-pkg.JwtSecret")

type claims struct {
	jwt.StandardClaims

	Account string            `json:"account,omitempty"` // 账号
	UserID  uint              `json:"uid"`               // 用户ID
	Name    string            `json:"name"`              // 姓名
	Roles   []string          `json:"roles"`             // 角色组
	IP      string            `json:"ip,omitempty"`      // 用户登录ID
	Extend  map[string]string `json:"extend,omitempty"`  // 扩展信息
}
