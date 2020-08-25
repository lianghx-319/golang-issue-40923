package table

import (
	. "golang-issues-40923/src/common"
	"strings"
)

type FilterData struct {
	DiffThreshold int32 // 单位微秒
	Types         int
	Name          string
	SortBy        int
	SortDesc      bool
}

func TypeOf(item *ItemData) int {
	if item.Treatment.Count == 0 {
		return OnlyDeleted
	}
	if item.Control.Count == 0 {
		return OnlyAdded
	}
	return OnlyChanged
}

func IsMatchFilter(item *ItemData, filter FilterData) bool {
	if filter.Types == 0 {
		return true
	}
	if AbsInt32(int32(item.Treatment.Total-item.Control.Total)) < filter.DiffThreshold {
		return false
	}
	if !ContainMethodStatus(filter.Types, TypeOf(item)) {
		return false
	}
	if len(filter.Name) > 0 && !strings.Contains(strings.ToLower(item.Signature), strings.ToLower(filter.Name)) {
		return false
	}
	return true
}

type Page struct {
	Index int
	Size  int
}
