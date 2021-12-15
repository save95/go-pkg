package dbcache

import (
	"context"

	"github.com/save95/go-pkg/model/pager"
)

type ICache interface {
	// Paginate 分页列表
	Paginate(ctx context.Context, opt pager.Option) (*Paginate, error)
	// First 按 id 查询数据
	First(ctx context.Context, id uint) (interface{}, error)
	// ClearAll 清理所有缓存
	ClearAll(ctx context.Context) error
	// ClearPaginate 清理所有分页查询缓存
	ClearPaginate(ctx context.Context) error
	// ClearFirst 清理指定数据缓存
	ClearFirst(ctx context.Context, id uint) error
}
