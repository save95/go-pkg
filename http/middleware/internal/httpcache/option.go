package httpcache

import (
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/save95/go-pkg/http/jwt"
	"github.com/save95/go-pkg/http/middleware/internal/httpcache/store"
	"github.com/save95/go-utils/sliceutil"
	"github.com/save95/xlog"
)

type Option func(*handler)

func WithRedisStoreBy(addr string, db uint) Option {
	return func(c *handler) {
		if len(addr) != 0 {
			c.store = store.NewRedisStore(redis.NewClient(&redis.Options{
				Network: "tcp",
				Addr:    addr,
				DB:      int(db),
			}))
		}
	}
}

func WithRedisStore(client *redis.Client) Option {
	return func(c *handler) {
		if client != nil {
			c.store = store.NewRedisStore(client)
		}
	}
}

func WithLogger(log xlog.XLogger) Option {
	return func(c *handler) {
		c.log = log
	}
}

func WithDebug(enabled bool) Option {
	return func(c *handler) {
		c.debug = enabled
	}
}

func WithJWTOption(opt *jwt.Option) Option {
	return func(c *handler) {
		c.jwtOption = opt
	}
}

func WithGlobalCacheDuration(d time.Duration) Option {
	return func(c *handler) {
		c.globalCacheDuration = d
	}
}

func WithGlobalSkipQueryFields(fields ...string) Option {
	return func(c *handler) {
		for _, field := range fields {
			c.globalSkipFields[field] = struct{}{}
		}
	}
}

func WithCacheKeyPrefix(str string) Option {
	return func(c *handler) {
		c.prefixKey = str
	}
}

func WithoutHeader(without bool) Option {
	return func(c *handler) {
		c.withoutHeader = without
	}
}

func WithRoutePolicy(route string, withToken bool, fields ...string) Option {
	return withRouteRule(route, withToken, 0, fields, nil)
}

func WithRouteRule(route string, withToken bool, duration time.Duration, fields ...string) Option {
	return withRouteRule(route, withToken, duration, fields, nil)
}

func WithRouteSkipFiledPolicy(route string, withToken bool, skipFields ...string) Option {
	return withRouteRule(route, withToken, 0, nil, skipFields)
}

func WithRouteSkipFiledRule(route string, withToken bool, duration time.Duration, skipFields ...string) Option {
	return withRouteRule(route, withToken, duration, nil, skipFields)
}

func withRouteRule(route string, withToken bool, duration time.Duration, fields, skipFields []string) Option {
	return func(c *handler) {
		// 先记录顺序
		c.routeList = append(c.routeList, route)
		c.routeList = sliceutil.UniqueString(c.routeList...)

		if c.routePolicies == nil {
			c.routePolicies = make(map[string]*ruleItem, 0)
		}

		rule, ok := c.routePolicies[route]
		if !ok {
			rule = &ruleItem{}
		}

		rule.withToken = withToken
		rule.duration = duration

		if fields != nil {
			if rule.fields == nil {
				rule.fields = make(map[string]struct{}, 0)
			}
			for _, field := range fields {
				rule.fields[field] = struct{}{}
			}
		}

		if skipFields != nil {
			if rule.skipFields == nil {
				rule.skipFields = make(map[string]struct{}, 0)
			}
			for _, field := range skipFields {
				rule.skipFields[field] = struct{}{}
			}
		}

		// 写入规则
		c.routePolicies[route] = rule
	}
}
