package middleware

import (
	"time"

	"github.com/save95/go-pkg/http/types"
)

// JWTOption jwt 配置参数
type JWTOption struct {
	// 角色转化函数，转换成角色 types.IRole
	RoleConvert types.ToRole
	// 过期自动刷新临界时长。零则表示不自动刷新
	RefreshDuration time.Duration
	// token 加密密钥。默认为 "go-pkg.JwtSecret"
	Secret []byte
}
