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
)

type dbCache struct {
	cacheManager *cache.Cache
	name         string
	autoRenew    bool
	expiration   time.Duration
	//log          xlog.XLogger
}

func NewDefault(name string, cacheManager *cache.Cache) *dbCache {
	return &dbCache{
		name:         name,
		cacheManager: cacheManager,
		autoRenew:    true,
		expiration:   5 * time.Minute,
	}
}

func (s *dbCache) WithAutoRenew(autoRenew bool) *dbCache {
	s.autoRenew = autoRenew
	return s
}

func (s *dbCache) WithExpiration(expiration time.Duration) *dbCache {
	if expiration == 0 {
		expiration = 5 * time.Minute
	}
	s.expiration = expiration
	return s
}

func (s *dbCache) Paginate(ctx context.Context, opt pager.Option,
	fun func(ctx context.Context, opt pager.Option) (interface{}, uint, error),
) (*Paginate, error) {
	kbs, _ := json.Marshal(opt)
	k := strings.ToLower(fmt.Sprintf("%x", md5.Sum(kbs)))
	key := fmt.Sprintf("%s:paginate:%s", s.name, k)

	// 登记缓存key，方便后续清理
	s.appendKey(ctx, key)

	var data Paginate
	cacheData, d, err := s.cacheManager.GetWithTTL(ctx, key)
	if nil == err {
		_ = json.Unmarshal([]byte(cacheData.(string)), &data)

		// 延长时效
		if s.autoRenew && d <= time.Minute {
			err = s.cacheManager.Set(ctx, key, cacheData, &store.Options{
				Expiration: s.expiration,
			})
		}

		return &data, nil
	}

	if redis.Nil != err {
		return nil, err
	}

	//s.log.Debugf("key[%s] no cache, query", key)
	records, total, err := fun(ctx, opt)
	if nil != err {
		return nil, err
	}

	slices, ok := sliceutil.ToAny(records)
	if !ok {
		slices = make([]interface{}, 0)
	}

	data = Paginate{
		Data:  slices,
		Total: total,
		Query: opt,
	}

	bs, _ := json.Marshal(data)

	err = s.cacheManager.Set(ctx, key, string(bs), &store.Options{
		Expiration: 5 * time.Minute,
	})
	if nil != err {
		return nil, xerror.Wrap(err, "cache store failed")
	}

	return &data, nil
}

func (s *dbCache) First(ctx context.Context, id uint,
	fun func(ctx context.Context, id uint) (interface{}, error),
) (interface{}, error) {
	if id == 0 {
		return nil, xerror.New("id error")
	}

	key := fmt.Sprintf("%s:first:%d", s.name, id)

	// 登记缓存key，方便后续清理
	s.appendKey(ctx, key)

	var data interface{}
	cacheData, d, err := s.cacheManager.GetWithTTL(ctx, key)
	if nil == err {
		_ = json.Unmarshal([]byte(cacheData.(string)), &data)

		// 延长时效
		if s.autoRenew && d <= time.Minute {
			err = s.cacheManager.Set(ctx, key, cacheData, &store.Options{
				Expiration: s.expiration,
			})
		}

		return &data, nil
	}

	if redis.Nil != err {
		return nil, err
	}

	//s.log.Debugf("key[%s] no cache, query", key)
	record, err := fun(ctx, id)
	if nil != err {
		return nil, err
	}

	data = record
	bs, _ := json.Marshal(data)

	err = s.cacheManager.Set(ctx, key, string(bs), &store.Options{
		Expiration: 5 * time.Minute,
	})
	if nil != err {
		return nil, xerror.Wrap(err, "cache store failed")
	}

	return &data, nil
}

func (s *dbCache) ClearAll(ctx context.Context) error {
	keys := s.keys(ctx)
	for _, key := range keys {
		if err := s.cacheManager.Delete(ctx, key); nil != err {
			return err
		}
	}

	return nil
}

func (s *dbCache) ClearPaginate(ctx context.Context) error {
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
	key := fmt.Sprintf("%s:first:%d", s.name, id)
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
