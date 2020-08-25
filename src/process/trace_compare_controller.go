package process

import (
	. "golang-issues-40923/src/common"
	"golang-issues-40923/src/parse"
	"golang-issues-40923/src/parse/base"
	. "golang-issues-40923/src/process/table"
	"sync"
)

type TraceCompareController struct {
	context         *traceContext
	tableController *TraceTableProcessor
	oldStackRoot    *MethodNode
	newStackRoot    *MethodNode
}

var mu sync.Mutex
var instance *TraceCompareController

func GetTraceCompareController() *TraceCompareController {
	mu.Lock()
	defer mu.Unlock()

	if instance == nil {
		ctx := &traceContext{
			controlParseCtx: &base.ParseContext{
				IdGen: NewIdGenerator(0, 1),
			},
			treatmentParseCtx: &base.ParseContext{
				IdGen: NewIdGenerator(1, 1),
			},
		}
		instance = &TraceCompareController{
			context:         ctx,
			tableController: nil,
		}
	}

	return instance
}

func (controller *TraceCompareController) Parse(chunk []byte, isEOF bool, traceGroup int) bool {
	controlParser := controller.context.controlParser
	treatmentParser := controller.context.treatmentParser
	needMore := true
	switch traceGroup {
	case GroupControl:
		if controlParser == nil {
			controlParser = parse.NewForwardingParser(controller.context.controlParseCtx)
			controller.context.controlParser = controlParser
		}
		needMore = controlParser.Parse(chunk, isEOF)
		if isEOF || !needMore {
			controller.oldStackRoot = controlParser.GetResult()
		}
	case GroupTreatment:
		if treatmentParser == nil {
			treatmentParser = parse.NewForwardingParser(controller.context.treatmentParseCtx)
			controller.context.treatmentParser = treatmentParser
		}
		needMore = treatmentParser.Parse(chunk, isEOF)
		if isEOF || !needMore {
			controller.newStackRoot = treatmentParser.GetResult()
		}
	}
	return needMore
}

func (controller *TraceCompareController) GetTableResult(filterMethod string, diffThreshold int, onlyChanged bool,
	onlyAdded bool, onlyDeleted bool, sortBy int, sortDesc bool, pageIndex int, pageSize int) string {
	controller.ensureTableController()
	return controller.tableController.GetTablePage(FilterData{
		DiffThreshold: int32(diffThreshold),
		Types:         GetMethodStatus(onlyChanged, onlyAdded, onlyDeleted),
		Name:          filterMethod,
		SortBy:        sortBy,
		SortDesc:      sortDesc,
	}, Page{
		Index: pageIndex,
		Size:  pageSize,
	})
}

func (controller *TraceCompareController) GetNewStackRoot() *MethodNode {
	return controller.newStackRoot
}

func (controller *TraceCompareController) GetOldStackRoot() *MethodNode {
	return controller.oldStackRoot
}

func (controller *TraceCompareController) ensureTableController() {
	if controller.tableController == nil {
		controller.tableController = NewTraceTableProcessor()
		controller.tableController.Process(controller.newStackRoot, controller.oldStackRoot)
	}
}

type traceContext struct {
	controlParseCtx   *base.ParseContext
	controlParser     *parse.ForwardingParser
	treatmentParseCtx *base.ParseContext
	treatmentParser   *parse.ForwardingParser
}
