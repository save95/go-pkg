package restful

type IResponse interface {
	// SetHeader 设置请求头
	SetHeader(key, value string) IResponse

	// Retrieve 查询单个资源的响应
	Retrieve(entity interface{})
	// TableWithPagination 表格分页响应
	TableWithPagination(resp *TableResponse)
	// ListWithPagination 分页列表的响应
	ListWithPagination(totalRow uint, entities interface{})
	// ListWithMoreFlag 查询列表的响应
	ListWithMoreFlag(hasMore bool, entities interface{})

	// Post 新增请求的响应
	Post(entity interface{})
	// Put 全量更新资源的响应
	Put(entity interface{})
	// Patch 部分更新资源的响应
	// 部分 cdn 服务商不支持 http patch 方法，如 阿里云
	Patch(entity interface{})
	// Delete 删除的响应
	Delete(err error)

	// WithMessage 通过 json 响应文本消息: {"message": "something..."}
	WithMessage(msg string)
	// WithBody 响应文本消息
	WithBody(body string)
	// WithError 响应错误消息(HttpStatus!=200)
	WithError(err error)
	// WithErrorData 响应错误消息(HttpStatus!=200)，并在 header 中返回错误数据
	WithErrorData(err error, data interface{})
}

type TableResponse struct {
	TotalRow uint                          // 分页的记录条数
	Columns  []string                      // 表格列
	RowKeys  []string                      // 表格行
	Items    []*TableResponseItem          // 表格行数据
	Extends  []*TableResponseRowExtendItem // 表格行扩展数据
}

type TableResponseItem struct {
	Column string      // 列
	RowKey string      // 行关键字
	Data   interface{} // 数据
}

type TableResponseRowExtendItem struct {
	RowKey string      // 行关键字
	Data   interface{} // 数据
}
