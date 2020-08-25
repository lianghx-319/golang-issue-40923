package table

import (
	. "golang-issues-40923/src/common"
	"sort"
	"strings"
)

type Snapshot struct {
	filter *FilterData
	data   []*ItemData
}

func (s *Snapshot) Match(filter *FilterData) bool {
	if filter == nil {
		return false
	}
	if filter.DiffThreshold != s.filter.DiffThreshold {
		return false
	}
	if filter.Types != s.filter.Types {
		return false
	}
	if strings.Compare(filter.Name, s.filter.Name) != 0 {
		return false
	}
	if filter.SortBy != s.filter.SortBy {
		return false
	}
	if filter.SortDesc != s.filter.SortDesc {
		return false
	}
	return true
}

func (s *Snapshot) GetPage(pageIndex int, pageSize int) (data []*ItemData, total int) {
	start := pageIndex * pageSize
	end := start + pageSize
	total = len(s.data)
	if end > total {
		end = total
	}
	if end > start {
		return s.data[start:end], total
	}
	return nil, total
}

func CreateSnapshot(filter FilterData, data map[string]*ItemData) *Snapshot {
	snap := Snapshot{}
	snap.filter = &filter
	list := make([]*ItemData, 0, 64)
	for _, item := range data {
		if IsMatchFilter(item, filter) {
			list = append(list, item)
		}
	}
	sort.SliceStable(list, func(i, j int) bool {
		var result bool
		switch filter.SortBy {
		case SortByDiffPartition:
			result = list[i].Diff.Partition > list[j].Diff.Partition
		case SortByDiffSelf:
			result = list[i].Diff.Self > list[j].Diff.Self
		case SortByDiffCount:
			result = list[i].Diff.Count > list[j].Diff.Count
		case SortByControlTotal:
			result = list[i].Control.Total > list[j].Control.Total
		case SortByControlSelf:
			result = list[i].Control.Self > list[j].Control.Self
		case SortByControlCount:
			result = list[i].Control.Count > list[j].Control.Count
		case SortByTreatmentTotal:
			result = list[i].Treatment.Total > list[j].Treatment.Total
		case SortByTreatmentSelf:
			result = list[i].Treatment.Self > list[j].Treatment.Self
		case SortByTreatmentCount:
			result = list[i].Treatment.Count > list[j].Treatment.Count
		default:
			result = list[i].Diff.Total > list[j].Diff.Total
		}
		return result && filter.SortDesc || !result && !filter.SortDesc
	})
	snap.data = list
	return &snap
}
