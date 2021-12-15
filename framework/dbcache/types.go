package dbcache

import (
	"github.com/save95/go-pkg/model/pager"
)

type Paginate struct {
	Query pager.Option
	Data  []interface{}
	Total uint
}
