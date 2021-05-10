package pager

// Sorted 排序顺序
type Sorted uint8

const (
	ASC  Sorted = iota // 正序
	DESC               // 倒序
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
