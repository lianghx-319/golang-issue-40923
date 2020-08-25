package table

import (
	"fmt"
	. "golang-issues-40923/src/common"
	"strconv"
)

type MethodData struct {
	Self     int64       `json:"self"`
	Total    int64       `json:"total"`
	Count    int32       `json:"count"`
	NodeList *MethodList `json:"-"`
}

func (data *MethodData) AddMethod(method *MethodNode) {
	data.Self += method.SelfTime
	data.Total += method.TotalTime
	data.Count++
	if data.NodeList == nil {
		data.NodeList = NewMethodList()
	}
	data.NodeList.Add(method)
}

type DiffData struct {
	Total     int64   `json:"total"`
	Partition float32 `json:"partition"`
	Self      int64   `json:"self"`
	Count     int32   `json:"count"`
}

func CalculateDiffData(treatment *MethodData, control *MethodData, totalDiff int64) *DiffData {
	diff := DiffData{}
	diff.Total = treatment.Total - control.Total
	if totalDiff == 0 {
		diff.Partition = 1
	} else {
		val, _ := strconv.ParseFloat(fmt.Sprintf("%.6f", float64(AbsInt64(diff.Total))/float64(AbsInt64(totalDiff))), 32)
		diff.Partition = float32(val)
	}
	diff.Self = treatment.Self - control.Self
	diff.Count = treatment.Count - control.Count
	return &diff
}

type ItemData struct {
	Signature string      `json:"name"`
	Diff      *DiffData   `json:"diff"`
	Treatment *MethodData `json:"treatment"`
	Control   *MethodData `json:"control"`
}

func NewItemData(signature string) *ItemData {
	return &ItemData{
		Signature: signature,
		Diff:      nil,
		Treatment: &MethodData{},
		Control:   &MethodData{},
	}
}
