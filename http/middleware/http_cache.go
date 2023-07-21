package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/save95/go-pkg/http/jwt"
	"github.com/save95/go-pkg/http/middleware/internal/httpcache"
	"github.com/save95/xlog"
)

// HttpCache http 响应缓存
//
// usage:
//   r.Use(middleware.HttpCache(
//   	middleware.WithHttpCacheDebug(),
//   	middleware.WithHttpCacheLogger(global.Log),
//   	middleware.WithHttpCacheJWTOption(global.JWTOption(false)),
//   	middleware.WithHttpCacheGlobalDuration(5*time.Minute),
//   	middleware.WithHttpCacheRedisStore(redis.NewClient(&redis.Options{
//   		Addr:     global.Config.HttpCache.Addr,
//   		Password: global.Config.HttpCache.Password,
//   		DB:       global.Config.HttpCache.DB,
//   	})),
//   	middleware.WithHttpCacheGlobalSkipFields("v"),
//   	middleware.WithHttpCacheRouteSkipFiledPolicy("/user/", true),
//   ))
func HttpCache(opts ...httpcache.Option) gin.HandlerFunc {
	return httpcache.New(opts...)
}

func WithHttpCacheRedisStoreBy(addr string, db uint) httpcache.Option {
	return httpcache.WithRedisStoreBy(addr, db)
}

func WithHttpCacheRedisStore(client *redis.Client) httpcache.Option {
	return httpcache.WithRedisStore(client)
}

func WithHttpCacheLogger(log xlog.XLogger) httpcache.Option {
	return httpcache.WithLogger(log)
}

func WithHttpCacheDebug() httpcache.Option {
	return httpcache.WithDebug(true)
}

func WithHttpCacheJWTOption(opt *jwt.Option) httpcache.Option {
	return httpcache.WithJWTOption(opt)
}

func WithHttpCacheGlobalDuration(d time.Duration) httpcache.Option {
	return httpcache.WithGlobalCacheDuration(d)
}

func WithHttpCacheGlobalSkipFields(fields ...string) httpcache.Option {
	return httpcache.WithGlobalSkipQueryFields(fields...)
}

func WithHttpCacheKeyPrefix(str string) httpcache.Option {
	return httpcache.WithCacheKeyPrefix(str)
}

func WithoutHttpCacheHeader(without bool) httpcache.Option {
	return httpcache.WithoutHeader(without)
}

func WithHttpCacheRoutePolicy(route string, withToken bool, fields ...string) httpcache.Option {
	return httpcache.WithRoutePolicy(route, withToken, fields...)
}

func WithHttpCacheRouteRule(route string, withToken bool, duration time.Duration, fields ...string) httpcache.Option {
	return httpcache.WithRouteRule(route, withToken, duration, fields...)
}

func WithHttpCacheRouteSkipFiledPolicy(route string, withToken bool, skipFields ...string) httpcache.Option {
	return httpcache.WithRouteSkipFiledPolicy(route, withToken, skipFields...)
}

func WithHttpCacheRouteSkipFiledRule(route string, withToken bool, duration time.Duration, fields ...string) httpcache.Option {
	return httpcache.WithRouteSkipFiledRule(route, withToken, duration, fields...)
}
