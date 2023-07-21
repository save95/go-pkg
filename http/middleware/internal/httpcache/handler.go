package httpcache

import (
	"bytes"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/save95/go-pkg/http/jwt"
	"github.com/save95/go-pkg/http/middleware/internal/httpcache/store"
	"github.com/save95/xerror"
	"github.com/save95/xlog"
)

type handler struct {
	debug               bool
	singleFlightTimeout time.Duration
	withoutHeader       bool
	prefixKey           string
	log                 xlog.XLogger

	store     store.ICacheStore
	jwtOption *jwt.Option

	globalCacheDuration time.Duration
	globalSkipFields    map[string]struct{} // 不用于计算缓存的 key

	routeList     []string             // 路由规则排序列表
	routePolicies map[string]*ruleItem // 路由特殊规则: urlPathRegrex => ruleItem
}

func New(opts ...Option) gin.HandlerFunc {
	f := &handler{
		globalCacheDuration: 5 * time.Minute,
		globalSkipFields:    make(map[string]struct{}, 0),
		routeList:           make([]string, 0),
		routePolicies:       make(map[string]*ruleItem, 0),
		singleFlightTimeout: 10 * time.Millisecond, // 100QPS
	}

	for _, opt := range opts {
		opt(f)
	}

	return func(c *gin.Context) {
		strategy, err := f.getCacheStrategy(c)
		if nil != err {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": "get http cache strategy failed: " + err.Error(),
			})
			return
		}

		if f.store == nil {
			c.Next()
			return
		}

		if !strategy.NeedCached {
			c.Next()
			return
		}

		cached, respCache, err := f.cached(c, strategy)
		if nil != err {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": "http cache handle failed: " + err.Error(),
			})
			return
		}

		if !cached {
			c.Next()
			return
		}

		c.Writer.WriteHeader(respCache.Status)

		if !f.withoutHeader {
			for key, values := range respCache.Header {
				for _, val := range values {
					c.Writer.Header().Set(key, val)
				}
			}
		}

		if _, err := c.Writer.Write(respCache.Data); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": "http cache handle failed: " + err.Error(),
			})
			return
		}

		// 跳出，不走后续的中间件
		c.Abort()
	}
}

func (h *handler) getCacheStrategy(ctx *gin.Context) (*strategy, error) {
	fullPath := ctx.FullPath()
	method := ctx.Request.Method

	// 只缓存 get/delete 成功的请求
	if method != http.MethodGet && method != http.MethodDelete {
		h.debugf("http method is not GET/DELETE, skip. url=[%s]%s", method, fullPath)
		return &strategy{
			NeedCached: false,
		}, nil
	}

	// 获取路由单独的缓存策略，如果存在多个，则以最后一个为准
	var rule *ruleItem
	for _, route := range h.routeList {
		if strings.Contains(fullPath, route) {
			rule = h.routePolicies[route]
		}
	}

	if rule == nil {
		h.debugf("not hit strategy, skip. url=%s", fullPath)
		return &strategy{
			NeedCached: false,
		}, nil
	}
	h.debugf("found cache strategy: rule=%s, url=%s", rule.String(), fullPath)

	qs := ctx.Request.URL.Query()
	params := url.Values{}
	for key := range qs {
		if len(rule.fields) == 0 {
			_, gok := h.globalSkipFields[key]
			_, ok := rule.skipFields[key]
			if !gok && !ok {
				params.Add(key, qs.Get(key))
			}
		} else {
			if _, ok := rule.fields[key]; !ok {
				params.Add(key, qs.Get(key))
			}
		}
	}

	var userID uint
	if rule.withToken {
		user, err := jwt.MustParseJWTUser(ctx, h.jwtOption)
		if nil != err {
			h.debugf("parse jwt user failed, err=%+v", err)
			return nil, xerror.Wrap(err, "parse jwt user failed")
		}
		userID = user.GetID()
	}

	cacheKey := ctx.Request.URL.Path + ":" + params.Encode()
	if userID > 0 {
		cacheKey += ":forUser:" + strconv.Itoa(int(userID))
	}
	h.debugf("get cache strategy input: qs=%s, key=%s", qs.Encode(), cacheKey)

	duration := h.globalCacheDuration
	if rule.duration > 0 {
		duration = rule.duration
	}

	return &strategy{
		NeedCached:    true,
		CacheKey:      cacheKey,
		CacheDuration: duration,
	}, nil
}

func (h *handler) debugf(format string, vals ...interface{}) {
	if !h.debug {
		return
	}

	if h.log != nil {
		h.log.Debugf("[httpcache] "+format, vals...)
		return
	}

	log.Printf("[httpcache] "+format+"\n", vals...)
}

func (h *handler) cached(c *gin.Context, strategy *strategy) (bool, *store.CachedResponse, error) {
	cacheKey := h.getCacheKey(strategy.CacheKey)

	data, err, _ := sf.Do(cacheKey, func() (interface{}, error) {
		// 限制 QPS = 1s/h.singleFlightTimeout
		if h.singleFlightTimeout > 0 {
			timer := time.AfterFunc(h.singleFlightTimeout, func() {
				sf.Forget(cacheKey)
			})
			defer timer.Stop()
		}

		// 先获取缓存
		respCache := store.CachedResponse{}
		err := h.store.Get(cacheKey, &respCache)
		if err == nil {
			h.debugf("hit cache, key=%s", cacheKey)
			return &respCache, nil
		}

		if err != store.ErrorCacheMiss {
			return nil, xerror.Wrapf(err, "get http cache failed, key=%s", cacheKey)
		}

		// 未获取到缓存，调用下一个请求链
		// 将自定义的响应写入器传递给 Gin 的下一个处理器，便于复制和缓存 response
		cacheWriter := newResponseWriter(c.Writer)
		c.Writer = cacheWriter
		c.Next()

		// 非成功请求，不进行缓存
		if c.Writer.Status() < 200 || c.Writer.Status() >= 300 {
			h.debugf("http request not success, skip. statusCode=%d", c.Writer.Status())
			return nil, nil
		}

		// 保存缓存
		resp := newCachedResponse(cacheWriter)
		if err := h.store.Set(cacheKey, resp, strategy.CacheDuration); err != nil {
			return nil, xerror.Wrapf(err, "set http cache failed, key=%s", cacheKey)
		}

		// 从请求链中获取的数据，直接跳过。防止响应重复
		h.debugf("not cache, save cache and redirect next, key=%s", cacheKey)
		return nil, nil
	})

	if nil != err {
		// 非 debug 模式，不阻塞
		if !h.debug {
			return false, nil, nil
		}
		return false, nil, err
	}

	// 不需要缓存
	if data == nil {
		return false, nil, nil
	}

	return true, data.(*store.CachedResponse), nil
}

func (h *handler) getCacheKey(key string) string {
	var bf bytes.Buffer
	bf.WriteString("httpCache:")
	if len(h.prefixKey) > 0 {
		bf.WriteString(h.prefixKey)
		bf.WriteString(":")
	}
	bf.WriteString(key)
	return bf.String()
}

func (h *handler) replyWithCache(c *gin.Context, respCache *store.CachedResponse) {
	c.Writer.WriteHeader(respCache.Status)

	if !h.withoutHeader {
		for key, values := range respCache.Header {
			for _, val := range values {
				c.Writer.Header().Set(key, val)
			}
		}
	}

	if _, err := c.Writer.Write(respCache.Data); err != nil {
		h.debugf("write response error: %s", err)
	}
}
