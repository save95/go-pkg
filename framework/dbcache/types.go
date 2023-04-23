package dbcache

import (
	"context"
	"time"

	"github.com/save95/go-pkg/model/pager"
)

type IDBCache interface {
	WithAutoRenew(autoRenew bool) IDBCache            // 缓存是否自动延长有效期
	WithExpiration(expiration time.Duration) IDBCache // 缓存有效期

	Paginate(
		ctx context.Context,
		opt pager.Option,
		fun func() (interface{}, uint, error),
	) (*PaginateResult, error) // 分页列表
	First(
		ctx context.Context,
		id uint,
		fun func() (interface{}, error),
	) (string, error) //  按 id 查询数据
	Remember(
		ctx context.Context,
		key string,
		fun func() (interface{}, error),
	) (interface{}, error)

	ClearAll(ctx context.Context) error            // 清理所有缓存
	ClearPaginate(ctx context.Context) error       // 清理所有分页查询缓存
	ClearFirst(ctx context.Context, id uint) error // 清理指定数据缓存
	Forget(ctx context.Context, key string) error  // 清理指定数据缓存
}

type PaginateResult struct {
	DataBytes []byte
	Total     uint
}
