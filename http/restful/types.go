package restful

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
