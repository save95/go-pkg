package types

import "github.com/save95/go-pkg/constant"

// SearchRequest 通用搜索请求
type SearchRequest struct {
	// 数据开始位置
	Start int64 `form:"start"`
	// 返回数据条数
	Limit int64 `form:"limit"`
	// 排序规则：sort=otc_type,-created_at,*custom
	// 以符号开头，可选符号：(+或空 正序）（- 倒序）（* 自定义复杂排序标识关键词）
	Sort string `form:"sort"`
}

// GetLimit 获得分页数据条数
func (sr SearchRequest) GetLimit() int64 {
	if sr.Limit <= 0 {
		return constant.DefaultRequestLimit
	}

	return sr.Limit
}
