package pager

import "github.com/save95/go-pkg/constant"

// Option 分页参数
type Option struct {
	Start   int
	Limit   int
	Filter  Filter
	Sorters []Sorter
}

func (po Option) GetLimit() int {
	if po.Limit <= 0 {
		return constant.DefaultPageSize
	}

	return po.Limit
}
