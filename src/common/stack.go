package common

import "container/list"

type MethodStack struct {
	data *list.List
}

func (stack *MethodStack) Peek() *MethodNode {
	e := stack.data.Back()
	if e != nil {
		method, _ := e.Value.(*MethodNode)
		return method
	}
	return nil
}

func (stack *MethodStack) Pop() *MethodNode {
	e := stack.data.Back()
	if e != nil {
		method, _ := e.Value.(*MethodNode)
		stack.data.Remove(e)
		return method
	}
	return nil
}

func (stack *MethodStack) Push(method *MethodNode) {
	stack.data.PushBack(method)
}

func (stack *MethodStack) NotEmpty() bool {
	return stack.data.Len() > 0
}

func (stack *MethodStack) Size() int {
	return stack.data.Len()
}

func NewMethodStack() *MethodStack {
	return &MethodStack{
		data: list.New(),
	}
}
