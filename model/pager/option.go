package pager

import "github.com/save95/go-pkg/constant"

// Option 分页参数
type Option struct {
	Start   int
	Limit   int
	Filter  Filter
	Sorters []Sorter

	Preloads []string
}

func (po Option) GetLimit() int {
	if po.Limit <= 0 {
		return constant.DefaultPageSize
	}

	return po.Limit
}

func (po Option) GetSorters() []Sorter {
	if po.Sorters == nil {
		return []Sorter{}
	}

	return po.Sorters
}

func (po Option) GetPreloads() []string {
	if po.Preloads == nil {
		return []string{}
	}

	return po.Preloads
}
