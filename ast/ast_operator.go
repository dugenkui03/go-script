package ast


// OperatorAssociativity  vm 中计算的时候用到结合性，所以结合性的值应该是可导出的
type OperatorAssociativity int

const (
	LeftAssociativity OperatorAssociativity = iota + 1
	RightAssociativity
)


// OperatorPriority 操作符优先级，值越大优先级越高
// 目前只有在语法解析节点使用，所以是不可导出的
// 运算符优先级通过ast结构有表现，所以理论上也是不用导出的
type OperatorPriority int

const (
	firstLevelOp  OperatorPriority = iota + 1 // +, -
	secondLevelOp                             // *, /, %

	// 最大优先级运算符+1
	highestLevelOpPlusOne
)

func (p OperatorPriority) getIncrement() OperatorPriority {
	if p.isHighestLevelOp() {
		// should not happen
	}

	return p + 1
}

// isHighestLevelOp 是否是最高优先级的运算符
func (p OperatorPriority) isHighestLevelOp() bool {
	return p == (highestLevelOpPlusOne - 1)
}

var unaryOperator = map[string]bool{
	"+": true, "-": true,
}

var firstOperator = map[string]bool{
	"+": true, "-": true,
}

var secondOperator = map[string]bool{
	"*": true, "/": true, "%": true,
}

var operatorByPriority = map[OperatorPriority]map[string]bool{
	firstLevelOp: firstOperator, secondLevelOp: secondOperator,
}
