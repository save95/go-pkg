package model

import "github.com/save95/go-pkg/constant"

// Sorted 排序顺序
type Sorted uint8

const (
	ASC  Sorted = iota // 正序
	DESC               // 倒序
)

// Sorter 排序器
type Sorter struct {
	Field  string
	Sorted Sorted
}

// Filter 过滤器
type Filter map[string]interface{}

// PagerOption 分页参数
type PagerOption struct {
	Start   int
	Limit   int
	Filter  Filter
	Sorters []Sorter
}

func (po PagerOption) GetLimit() int {
	if po.Limit <= 0 {
		return constant.DefaultPageSize
	}

	return po.Limit
}
