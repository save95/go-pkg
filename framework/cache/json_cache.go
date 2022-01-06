package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/eko/gocache/v2/cache"
	"github.com/eko/gocache/v2/store"
	"github.com/go-redis/redis/v8"
	"github.com/save95/go-utils/strutil"
	"github.com/save95/xerror"
)

type jsonCache struct {
	cacheManager *cache.Cache
	name         string
}

func NewJson(nameSpace string, cacheManager *cache.Cache) ICache {
	return &jsonCache{
		name:         nameSpace,
		cacheManager: cacheManager,
	}
}

func (c *jsonCache) getKey(key string) string {
	return fmt.Sprintf("%s:%s", strutil.Camel(c.name), key)
}

func (c *jsonCache) Get(ctx context.Context, key string) (jsonStr string, ttl time.Duration, err error) {
	if c.cacheManager == nil {
		return "", 0, xerror.New("cache manager no init")
	}

	key = c.getKey(key)
	cacheData, d, err := c.cacheManager.GetWithTTL(ctx, key)
	if nil == err {
		return cacheData.(string), d, nil
	}

	if redis.Nil == err {
		return "", 0, nil
	}

	return "", 0, err
}

func (c *jsonCache) Pull(ctx context.Context, key string) (jsonStr string, err error) {
	if c.cacheManager == nil {
		return "", xerror.New("cache manager no init")
	}

	key = c.getKey(key)
	cacheData, err := c.cacheManager.Get(ctx, key)
	if nil == err {
		if err := c.cacheManager.Delete(ctx, key); nil != err {
			return "", err
		}

		return cacheData.(string), nil
	}

	if redis.Nil == err {
		return "", nil
	}

	return "", err
}

func (c *jsonCache) Set(ctx context.Context, key string, val interface{}, expire time.Duration) error {
	if c.cacheManager == nil {
		return xerror.New("cache manager no init")
	}

	// 禁止设置永久缓存
	if expire == 0 {
		expire = 5 * time.Minute
	}

	bs, err := json.Marshal(val)
	if err != err {
		return xerror.Wrap(err, "cache val marshal failed")
	}

	key = c.getKey(key)
	return c.cacheManager.Set(ctx, key, string(bs), &store.Options{
		Expiration: expire,
	})
}

func (c *jsonCache) Remember(
	ctx context.Context,
	key string,
	expire time.Duration,
	fun func(ctx context.Context) (interface{}, error),
) (jsonStr string, err error) {
	if c.cacheManager == nil {
		return "", xerror.New("cache manager no init")
	}

	key = c.getKey(key)
	cd, err := c.cacheManager.Get(ctx, key)
	if nil == err {
		return cd.(string), nil
	}

	v, err, _ := single.Do(key, func() (interface{}, error) {
		return fun(ctx)
	})
	if nil != err {
		return "", err
	}

	bs, err := json.Marshal(v)
	if err != nil {
		return "", err
	}

	// 禁止设置永久缓存
	if expire == 0 {
		expire = 5 * time.Minute
	}

	// 保存缓存
	err = c.cacheManager.Set(ctx, key, string(bs), &store.Options{
		Expiration: expire,
	})
	if nil != err {
		return "", err
	}

	return string(bs), nil
}

func (c *jsonCache) Clear(ctx context.Context, key string) error {
	if c.cacheManager == nil {
		return xerror.New("cache manager no init")
	}

	key = c.getKey(key)
	return c.cacheManager.Delete(ctx, key)
}
