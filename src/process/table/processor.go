package table

import (
	"encoding/json"
	. "golang-issues-40923/src/common"
)

type TraceTableProcessor struct {
	data          map[string]*ItemData
	totalDiffTime int64
	snapshot      *Snapshot
}

func (processor *TraceTableProcessor) getOrCreateItem(method *MethodNode) *ItemData {
	item := processor.data[method.Signature]
	if item == nil {
		item = NewItemData(method.Signature)
		processor.data[method.Signature] = item
	}
	return item
}

func (processor *TraceTableProcessor) Process(treatmentRoot *MethodNode, controlRoot *MethodNode) {
	if treatmentRoot != nil {
		treatmentRoot.ForEach(func(method *MethodNode) {
			if method == treatmentRoot {
				return
			}
			item := processor.getOrCreateItem(method)
			item.Treatment.AddMethod(method)
		})
	}
	if controlRoot != nil {
		controlRoot.ForEach(func(method *MethodNode) {
			if method == controlRoot {
				return
			}
			item := processor.getOrCreateItem(method)
			item.Control.AddMethod(method)
		})
	}
	if treatmentRoot != nil && controlRoot != nil {
		processor.totalDiffTime = treatmentRoot.TotalTime - controlRoot.TotalTime
	}
	for _, item := range processor.data {
		item.Diff = CalculateDiffData(item.Treatment, item.Control, processor.totalDiffTime)
	}
}

func (processor *TraceTableProcessor) GetTablePage(filterData FilterData, page Page) string {
	result := PageResult{}
	filterData = FilterData{
		DiffThreshold: 5,
		Types:         0,
		Name:          "",
		SortBy:        0,
		SortDesc:      false,
	}
	page = Page{
		Index: 0,
		Size:  1000000000,
	}
	if processor.snapshot == nil || !processor.snapshot.Match(&filterData) {
		processor.snapshot = CreateSnapshot(filterData, processor.data)
	}
	data, total := processor.snapshot.GetPage(page.Index, page.Size)
	result.List = data
	result.TotalSize = total
	result.PageIndex = page.Index
	result.PageSize = len(data)
	result.TotalDiffTime = processor.totalDiffTime
	resultJson, err := json.Marshal(result)
	if err != nil {
		panic(err)
	}
	return string(resultJson)
}

func NewTraceTableProcessor() *TraceTableProcessor {
	return &TraceTableProcessor{
		totalDiffTime: 0,
		data:          make(map[string]*ItemData),
	}
}
