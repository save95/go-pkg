package jwt

import (
	"time"

	"github.com/save95/go-pkg/http/types"
)

// Option jwt 配置参数
type Option struct {
	// 角色转化函数，转换成角色 IRole
	RoleConvert types.ToRole

	// 过期自动刷新临界时长。零则表示不自动刷新
	RefreshDuration time.Duration

	// token 加密密钥。默认为 "go-pkg.JwtSecret"
	Secret []byte

	// 是否开启静默模式。true-开启：鉴权失败，不注入用户信息；false-关闭。鉴权失败阻断，并抛出错误
	SilentMode bool
}

// WithSilent 设置是否开启静默模式
func (o *Option) WithSilent(enable bool) *Option {
	o.SilentMode = enable

	return o
}
