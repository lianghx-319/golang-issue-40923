package common

import (
	"container/list"
	"math"
)

type MethodNode struct {
	//方法基础字段
	Id             uint64
	Signature      string      //方法的签名 (boolean java.lang.reflect.Modifier.isPublic(int))
	EnterTimestamp uint64      //进入该方法的时间戳，单位微秒
	ExitTimestamp  uint64      //退出该方法的时间戳，单位微秒
	SelfTime       int64       //该方法自身耗时
	TotalTime      int64       //该方法总耗时
	IsTop          bool        //是否为顶层方法root
	ParentNode     *MethodNode //父结点
	Calls          *MethodList //内部调用的子方法列表
	//topdown拓展字段
	DelCalls         *MethodList //内部已删除的子方法列表
	MethodTag        int         //方法标签
	TotalChangedTime int64       //方法总变化时间
	SelfChangedTime  int64       //方法自身变化时间
	Suitability      int         //方法匹配度
}

func (method *MethodNode) SetEnter(timestamp uint64) {
	method.EnterTimestamp = timestamp
}

func (method *MethodNode) SetExit(timestamp uint64) {
	method.ExitTimestamp = timestamp
	method.TotalTime = int64(method.ExitTimestamp - method.EnterTimestamp)
	selfTime := method.TotalTime
	method.Calls.ForEach(func(child *MethodNode) {
		selfTime -= child.TotalTime
	})
	method.SelfTime = selfTime
}

func (method *MethodNode) AddCall(call *MethodNode) {
	method.Calls.Add(call)
}

func (method *MethodNode) AddDelCall(call *MethodNode) {
	method.DelCalls.Add(call)
}

func (method *MethodNode) GetAllCallsCount() int {
	return method.Calls.Len() + method.DelCalls.Len()
}

func (method *MethodNode) GetAllCalls() *MethodList {
	allCalls := NewMethodList()
	allCalls.Data.PushFrontList(method.Calls.Data)
	allCalls.Data.PushBackList(method.DelCalls.Data)
	return allCalls
}

func (method *MethodNode) GetAllCallsSortByCallTime() *MethodList {
	allCalls := method.GetAllCalls()
	allCalls.Data = sort(allCalls.Data, func(method1 *MethodNode, method2 *MethodNode) bool {
		return method1.EnterTimestamp > method2.EnterTimestamp
	})
	return allCalls
}

func (method *MethodNode) GetAllCallsSortByContribution() *MethodList {
	allCalls := method.GetAllCalls()
	allCalls.Data = sort(allCalls.Data, func(method1 *MethodNode, method2 *MethodNode) bool {
		return math.Abs(float64(method1.TotalChangedTime)) < math.Abs(float64(method2.TotalChangedTime))
	})
	return allCalls
}

func (method *MethodNode) ForEach(action func(method *MethodNode)) {
	stack := NewMethodStack()
	stack.Push(method)
	for stack.NotEmpty() {
		node := stack.Pop()
		action(node)
		node.Calls.ForEachReverse(func(child *MethodNode) {
			stack.Push(child)
		})
	}
}

func (method *MethodNode) AddMethodTag(methodTag int) {
	method.MethodTag |= methodTag
}

func (method *MethodNode) ContainMethodTag(methodTag int) bool {
	return method.MethodTag&methodTag == methodTag
}

func (method *MethodNode) IsMethodTagEmpty() bool {
	return method.MethodTag == 0
}

func NewMethodNode(id uint64, signature string) *MethodNode {
	return &MethodNode{
		Id:        id,
		Signature: signature,
		Calls:     NewMethodList(),
		DelCalls:  NewMethodList(),
	}
}

type MethodList struct {
	Data *list.List
}

func (list *MethodList) ForEach(action func(method *MethodNode)) {
	for e := list.Data.Front(); e != nil; e = e.Next() {
		method, _ := e.Value.(*MethodNode)
		action(method)
	}
}

func (list *MethodList) ForEachReverse(action func(method *MethodNode)) {
	for e := list.Data.Back(); e != nil; e = e.Prev() {
		method, _ := e.Value.(*MethodNode)
		action(method)
	}
}

func (list *MethodList) Add(method *MethodNode) {
	list.Data.PushBack(method)
}

func (list *MethodList) Len() int {
	return list.Data.Len()
}

func (list *MethodList) First() *MethodNode {
	front := list.Data.Front()
	if front == nil {
		return nil
	}
	m, _ := front.Value.(*MethodNode)
	return m
}

func (list *MethodList) Last() *MethodNode {
	back := list.Data.Back()
	if back == nil {
		return nil
	}
	m, _ := back.Value.(*MethodNode)
	return m
}

func (list *MethodList) Copy() *MethodList {
	newList := NewMethodList()
	if list != nil {
		list.ForEach(func(method *MethodNode) {
			newList.Add(method)
		})
	}
	return newList
}

func NewMethodList() *MethodList {
	return &MethodList{
		Data: list.New(),
	}
}

func sort(oldList *list.List, compare func(method1 *MethodNode, method2 *MethodNode) bool) (newList *list.List) {
	newList = list.New()
	for v := oldList.Front(); v != nil; v = v.Next() {
		node := newList.Front()
		for nil != node {
			if compare(node.Value.(*MethodNode), v.Value.(*MethodNode)) {
				newList.InsertBefore(v.Value.(*MethodNode), node)
				break
			}
			node = node.Next()
		}
		if node == nil {
			newList.PushBack(v.Value.(*MethodNode))
		}
	}
	return newList
}
