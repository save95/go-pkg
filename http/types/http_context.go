package types

import (
	"time"

	"github.com/save95/go-pkg/constant"

	"github.com/gin-gonic/gin"
	"github.com/save95/xlog"
)

type HttpContext struct {
	traceId      string       // 请求唯一标识
	version      ApiVersion   // 版本号
	bodyProperty BodyProperty // 响应正文属性
	user         *User        // 用户信息

	storage map[string]interface{} // 存储变量

	logger xlog.XLogger
}

func (sc *HttpContext) TraceId() string {
	return sc.traceId
}

func (sc *HttpContext) User() *User {
	return sc.user
}

func (sc *HttpContext) Logger() xlog.XLogger {
	return sc.logger
}

// HasRole 判断用户角色
func (sc *HttpContext) HasRole(roles []IRole) bool {
	if sc.user == nil {
		return false
	}

	for i := range roles {
		if sc.IsRole(roles[i]) {
			return true
		}
	}

	return false
}

// IsRole 判断用户角色
func (sc *HttpContext) IsRole(role IRole) bool {
	if sc.user == nil {
		return false
	}

	for i := range sc.user.Roles {
		if sc.user.Roles[i] == role {
			return true
		}
	}

	return false
}

func (sc *HttpContext) Deadline() (deadline time.Time, ok bool) {
	return
}

func (sc *HttpContext) Done() <-chan struct{} {
	return nil
}

func (sc *HttpContext) Err() error {
	return nil
}

func (sc *HttpContext) Value(key interface{}) interface{} {
	if keyAsString, ok := key.(string); ok {
		val, _ := sc.storage[keyAsString]
		return val
	}

	return nil
}

// Set 设置参数
// 会根据 value 的类型，自动设置对应属性的值，目前支持： ApiVersion, BodyProperty, types.User, xlog.XLogger
func (sc *HttpContext) Set(key string, value interface{}) {
	if len(key) == 0 {
		return
	}

	// API 版本号
	if av, ok := value.(ApiVersion); ok && av.Verify() {
		sc.version = av
		return
	}

	// 响应正文属性
	if bp, ok := value.(BodyProperty); ok && bp.Verify() {
		sc.bodyProperty = bp
		return
	}

	// 用户信息
	if v, ok := value.(User); ok {
		sc.user = &v
		return
	}

	// 日志处理器
	if v, ok := value.(xlog.XLogger); ok {
		sc.logger = v
	}

	// 延迟初始化
	if sc.storage == nil {
		sc.storage = make(map[string]interface{})
	}

	sc.storage[key] = value
}

// StorageTo 将已变更的数据，存储到 gin 上下文中，继续传输
func (sc *HttpContext) StorageTo(ctx *gin.Context) bool {
	ctx.Set(constant.HttpCustomContextKey, sc)
	return true
}
