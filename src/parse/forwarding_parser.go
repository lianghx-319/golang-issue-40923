package parse

import (
	. "golang-issues-40923/src/common"
	"golang-issues-40923/src/parse/base"
	"golang-issues-40923/src/parse/nanoscope"
)

type ForwardingParser struct {
	realParser base.Parser
	context    *base.ParseContext
}

func (parser *ForwardingParser) Parse(chunk []byte, isEOF bool) bool {
	if parser.realParser == nil {
		parser.realParser = nanoscope.NewParser(parser.context)
	}
	needMore := parser.realParser.Parse(chunk)
	if isEOF {
		parser.realParser.Finish()
	}
	return needMore
}

func (parser *ForwardingParser) GetResult() *MethodNode {
	return parser.realParser.Root()
}

func NewForwardingParser(ctx *base.ParseContext) *ForwardingParser {
	return &ForwardingParser{
		realParser: nil,
		context:    ctx,
	}
}
