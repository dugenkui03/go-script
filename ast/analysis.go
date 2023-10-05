package ast

// Equal 判断两个表达式是否相等
func Equal(e1, e2 Expression) bool {
	panic("没有实现")
}

// GetVariable 获取表达式使用的变量名称列表
func GetVariable(exp Expression) []string {
	panic("没有实现")
}

// GetFuncNames 获取表达式的使用的函数名称列表
func GetFuncNames(exp Expression) []string {
	panic("没有实现")
}
