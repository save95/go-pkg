package pager

import (
	"strings"

	"github.com/save95/go-utils/strutil"
)

// Sorted 排序顺序
type Sorted uint8

const (
	ASC    Sorted = iota // 正序
	DESC                 // 倒序
	Custom               // 自定义
)

func (s Sorted) String() string {
	switch s {
	case ASC:
		return "ASC"
	case DESC:
		return "DESC"
	default:
		return "ASC"
	}
}

// Sorter 排序器
type Sorter struct {
	Field  string
	Sorted Sorted
}

// ParseSorts 解析排序规则
// 以符号开头，可选符号：(+或空 正序）（- 倒序）（* 自定义复杂排序标识关键词）
// 多个排序规则按英文逗号隔开
func ParseSorts(sort string) []Sorter {
	sorters := make([]Sorter, 0)
	if len(sort) == 0 {
		return sorters
	}

	sorts := strings.Split(sort, ",")
	for _, s := range sorts {
		// query string 中的 + 在 net/url/url.go 中会被解析成空格
		s = strings.TrimSpace(s)
		switch s[:1] {
		case "*":
			sorters = append(sorters, Sorter{
				Field:  strutil.Snake(s[1:]),
				Sorted: Custom,
			})
		case "-":
			sorters = append(sorters, Sorter{
				Field:  strutil.Snake(s[1:]),
				Sorted: DESC,
			})
		case "+":
			sorters = append(sorters, Sorter{
				Field:  strutil.Snake(s[1:]),
				Sorted: ASC,
			})
		default:
			sorters = append(sorters, Sorter{
				Field:  strutil.Snake(s),
				Sorted: ASC,
			})
		}
	}

	return sorters
}
