package dbcache

import (
	"context"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/eko/gocache/v2/cache"
	"github.com/eko/gocache/v2/store"
	"github.com/go-redis/redis/v8"
	"github.com/save95/go-pkg/model/pager"
	"github.com/save95/go-utils/sliceutil"
	"github.com/save95/xerror"
	"golang.org/x/sync/singleflight"
)

var single singleflight.Group

type dbCache struct {
	cacheManager *cache.Cache
	name         string
	autoRenew    bool // 自动延长缓存有效期
	expiration   time.Duration
	//log          xlog.XLogger
}

func NewDefault(name string, cacheManager *cache.Cache) IDBCache {
	return &dbCache{
		name:         name,
		cacheManager: cacheManager,
		autoRenew:    true,
		expiration:   5 * time.Minute,
	}
}

func (s *dbCache) WithAutoRenew(autoRenew bool) IDBCache {
	s.autoRenew = autoRenew
	return s
}

func (s *dbCache) WithExpiration(expiration time.Duration) IDBCache {
	if expiration == 0 {
		expiration = 5 * time.Minute
	}
	s.expiration = expiration
	return s
}

func (s *dbCache) Paginate(ctx context.Context, opt pager.Option,
	fun func() (interface{}, uint, error),
) (*PaginateResult, error) {
	kbs, _ := json.Marshal(opt)
	k := strings.ToLower(fmt.Sprintf("%x", md5.Sum(kbs)))
	key := fmt.Sprintf("%s:paginate:%s", s.name, k)

	jsonStr, err := s.getOrQuery(ctx, key, func() (interface{}, error) {
		records, total, err := fun()
		if nil != err {
			return nil, err
		}

		slices, ok := sliceutil.ToAny(records)
		if !ok {
			slices = make([]interface{}, 0)
		}

		return &PaginateResult{
			Data:  slices,
			Total: total,
			Query: opt,
		}, nil
	})
	if nil != err {
		return nil, err
	}

	var data PaginateResult
	if err := json.Unmarshal([]byte(jsonStr), &data); nil != err {
		return nil, xerror.Wrap(err, "data unmarshal failed")
	}

	return &data, nil
}

func (s *dbCache) First(ctx context.Context, id uint,
	fun func() (interface{}, error),
) (interface{}, error) {
	if id == 0 {
		return nil, xerror.New("id error")
	}

	key := fmt.Sprintf("%s:first:%d", s.name, id)
	jsonStr, err := s.getOrQuery(ctx, key, fun)
	if nil != err {
		return nil, err
	}

	var data interface{}
	if err := json.Unmarshal([]byte(jsonStr), &data); nil != err {
		return nil, xerror.Wrap(err, "data unmarshal failed")
	}

	return data, nil
}

func (s *dbCache) Remember(ctx context.Context, key string,
	fun func() (interface{}, error),
) (interface{}, error) {
	key = fmt.Sprintf("%s:remember:%s", s.name, key)

	jsonStr, err := s.getOrQuery(ctx, key, fun)
	if nil != err {
		return nil, err
	}

	var data interface{}
	if err := json.Unmarshal([]byte(jsonStr), &data); nil != err {
		return nil, xerror.Wrap(err, "data unmarshal failed")
	}

	return data, nil
}

func (s *dbCache) getOrQuery(ctx context.Context, key string,
	fun func() (interface{}, error),
) (string, error) {
	if s.cacheManager == nil {
		return "", xerror.New("cache manager no init")
	}

	// 登记缓存key，方便后续清理
	s.appendKey(ctx, key)

	cacheData, d, err := s.cacheManager.GetWithTTL(ctx, key)
	if nil == err {
		// 延长时效
		if s.autoRenew && d <= time.Minute {
			err = s.cacheManager.Set(ctx, key, cacheData, &store.Options{
				Expiration: s.expiration,
			})
		}

		return cacheData.(string), nil
	}

	if redis.Nil != err {
		return "", err
	}

	v, err, _ := single.Do(key, func() (interface{}, error) {
		//s.log.Debugf("key[%s] no cache, query", key)
		record, err := fun()
		if nil != err {
			return nil, err
		}

		bs, err := json.Marshal(record)
		if nil != err {
			return nil, err
		}

		err = s.cacheManager.Set(ctx, key, string(bs), &store.Options{
			Expiration: s.expiration,
		})
		if nil != err {
			return nil, err
		}

		return string(bs), nil
	})
	if nil != err {
		return "", xerror.Wrap(err, "cache store failed")
	}

	return v.(string), nil
}

func (s *dbCache) ClearAll(ctx context.Context) error {
	if s.cacheManager == nil {
		return xerror.New("cache manager no init")
	}

	keys := s.keys(ctx)
	for _, key := range keys {
		if err := s.cacheManager.Delete(ctx, key); nil != err {
			return err
		}
	}

	return nil
}

func (s *dbCache) ClearPaginate(ctx context.Context) error {
	if s.cacheManager == nil {
		return xerror.New("cache manager no init")
	}

	keys := s.keys(ctx)
	for _, key := range keys {
		if strings.Contains(key, fmt.Sprintf("%s:paginate:", s.name)) {
			if err := s.cacheManager.Delete(ctx, key); nil != err {
				return err
			}
		}
	}

	return nil
}

func (s *dbCache) ClearFirst(ctx context.Context, id uint) error {
	if s.cacheManager == nil {
		return xerror.New("cache manager no init")
	}

	key := fmt.Sprintf("%s:first:%d", s.name, id)
	return s.cacheManager.Delete(ctx, key)
}

func (s *dbCache) Forget(ctx context.Context, key string) error {
	if s.cacheManager == nil {
		return xerror.New("cache manager no init")
	}

	key = fmt.Sprintf("%s:remember:%s", s.name, key)
	return s.cacheManager.Delete(ctx, key)
}

func (s *dbCache) appendKey(ctx context.Context, key string) {
	allKey := fmt.Sprintf("%s:cacheKeys", s.name)

	keys := s.keys(ctx)
	keys = append(keys, key)

	// 去重
	sets := make(map[string]struct{}, 0)
	for _, s2 := range keys {
		sets[s2] = struct{}{}
	}
	uniques := make([]string, 0)
	for s2 := range sets {
		uniques = append(uniques, s2)
	}

	bs, _ := json.Marshal(uniques)
	_ = s.cacheManager.Set(ctx, allKey, string(bs), &store.Options{
		Expiration: 5 * time.Minute,
	})
}

func (s *dbCache) keys(ctx context.Context) []string {
	allKey := fmt.Sprintf("%s:cacheKeys", s.name)
	cacheData, err := s.cacheManager.Get(ctx, allKey)
	if nil != err {
		cacheData = "[]"
	}

	var keys []string
	_ = json.Unmarshal([]byte(cacheData.(string)), &keys)
	if nil == keys {
		keys = make([]string, 0)
	}

	return keys
}
