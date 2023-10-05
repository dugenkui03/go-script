package vm

type Stack struct {
	value []interface{}
}

func NewStack() *Stack {
	return &Stack{value: make([]interface{}, 0)}
}

func (stack *Stack) push(ele interface{}) {
	stack.value = append(stack.value, ele)
}

func (stack *Stack) pop() interface{} {
	if stack == nil || len(stack.value) == 0 {
		return nil
	}

	top := stack.value[len(stack.value)-1]         // 获取栈顶元素
	stack.value = stack.value[:len(stack.value)-1] // 弹出栈顶元素
	return top
}
