package pager

import "testing"

func TestParseSorts(t *testing.T) {
	str := "createdAt,-id,*owner"
	t.Log(ParseSorts(str))
	t.Log(ParseSorts(""))
}
