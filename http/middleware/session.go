package middleware

import (
	"net/http"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

type SessionOption struct {
	Path     string
	Domain   string
	MaxAge   time.Duration
	Secure   bool
	HttpOnly bool
	SameSite http.SameSite
}

// Session 校验 session
// keyPairs cookie 键名
// secret cookie 存储加密密钥
func Session(keyPairs, secret string, opt SessionOption) gin.HandlerFunc {
	if len(secret) == 0 {
		secret = "go-pkg"
	}

	var store sessions.Store
	store = cookie.NewStore([]byte(secret))
	store.Options(sessions.Options{
		Path:     opt.Path,
		Domain:   opt.Domain,
		MaxAge:   int(opt.MaxAge.Seconds()), //seconds
		Secure:   opt.Secure,
		HttpOnly: opt.HttpOnly,
		SameSite: opt.SameSite,
	})

	return sessions.Sessions(keyPairs, store)
}

// SessionWithStore 校验 session
// keyPairs cookie 键名
func SessionWithStore(keyPairs string, store sessions.Store, opt SessionOption) gin.HandlerFunc {
	store.Options(sessions.Options{
		Path:     opt.Path,
		Domain:   opt.Domain,
		MaxAge:   int(opt.MaxAge.Seconds()), //seconds
		Secure:   opt.Secure,
		HttpOnly: opt.HttpOnly,
		SameSite: opt.SameSite,
	})

	return sessions.Sessions(keyPairs, store)
}
