package base

import . "golang-issues-40923/src/common"

type Parser interface {
	Parse(data []byte) bool //true:继续读取，false:可以停止读取数据
	Finish()
	Root() *MethodNode
}
