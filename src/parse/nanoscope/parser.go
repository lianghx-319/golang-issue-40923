package nanoscope

import (
	. "golang-issues-40923/src/common"
	"golang-issues-40923/src/parse/base"
	"strconv"
	"strings"
)

type Parser struct {
	context      *base.ParseContext
	root         *MethodNode
	stack        *MethodStack
	buffer       *base.MultiChunkBuffer
	limitReached bool
	tsBase       uint64 // 单位微秒
	tsRealTime   uint64 // 单位微秒
}

func (parser *Parser) Parse(data []byte) bool {
	parser.buffer.Write(data)
	parser.parseAllLines()
	return !parser.limitReached
}

func (parser *Parser) Finish() {
	parser.buffer.WriteEOF()
	parser.parseAllLines() // 读取剩余数据
	parser.parseEnd()
}

func (parser *Parser) parseEnd() {
	first := parser.root.Calls.First()
	if first != nil {
		parser.root.SetEnter(first.EnterTimestamp)
	}
	last := parser.root.Calls.Last()
	if last != nil {
		parser.root.SetExit(last.ExitTimestamp)
	}
}

func (parser *Parser) parseAllLines() {
	for {
		lineBytes := parser.buffer.ReadLine()
		if lineBytes == nil {
			break
		}
		line := string(lineBytes)
		parser.parseLine(line)
	}
}

func (parser *Parser) parseLine(line string) {
	if strings.HasPrefix(line, "#") {
		parser.parseExtLine(line)
	} else {
		parser.parseMethodLine(line)
	}
}

func (parser *Parser) parseExtLine(line string) {
	content := strings.TrimLeft(line, "# ")
	strBlocks := strings.Split(content, ":")
	if len(strBlocks) < 2 || len(strBlocks[0]) == 0 {
		return
	}
	key := strBlocks[0]
	value := strBlocks[1]
	if strings.Compare(key, "now_in_seconds") == 0 {
		s, err := strconv.ParseFloat(value, 32)
		if err != nil {
			return
		}
		parser.tsBase = uint64(s * 1000000) // 秒（浮点型）转为微秒（整型）
	} else if strings.Compare(key, "realtime_in_ms") == 0 {
		ms, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return
		}
		parser.tsRealTime = ms * 1000
	}
}

func (parser *Parser) parseMethodLine(line string) {
	strBlocks := strings.Split(line, ":")
	if len(strBlocks) < 2 || len(strBlocks[0]) == 0 || len(strBlocks[1]) == 0 {
		return
	}
	ns, err := strconv.ParseUint(strBlocks[0], 10, 64)
	if err != nil {
		return
	}
	timestamp := ns/1000 - parser.tsBase + parser.tsRealTime
	if strings.EqualFold(strBlocks[1], "POP") {
		method := parser.stack.Pop()
		if method != nil {
			method.SetExit(timestamp)
			pMethod := parser.stack.Peek()
			if pMethod == nil {
				limitTimestamp := parser.context.LimitTimestamp
				if limitTimestamp > 0 && timestamp > limitTimestamp {
					parser.limitReached = true
					parser.parseEnd()
					return
				}
				method.ParentNode = parser.root
				parser.root.AddCall(method)
			} else {
				method.ParentNode = pMethod
				pMethod.AddCall(method)
			}
		}
	} else {
		method := NewMethodNode(parser.context.IdGen.GenId(), strBlocks[1])
		method.SetEnter(timestamp)
		parser.stack.Push(method)
	}
}

func (parser *Parser) Root() *MethodNode {
	return parser.root
}

func NewParser(ctx *base.ParseContext) *Parser {
	return &Parser{
		context:      ctx,
		root:         NewMethodNode(ctx.IdGen.GenId(), "root"),
		stack:        NewMethodStack(),
		buffer:       base.NewMultiChunkBuffer(),
		limitReached: false,
	}
}
