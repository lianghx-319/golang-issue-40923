package base

import "golang-issues-40923/src/common"

type ParseContext struct {
	LimitTimestamp uint64
	IdGen          *common.IdGenerator
}
