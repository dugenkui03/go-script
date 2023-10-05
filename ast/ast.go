// Package ast
// 抽象语法数相关定义。通过语法推导ast结构定义的原则见笔记 page-4, page-5。简要描述如下
//
// 有文法
// A -> B
// B -> C
//
//	| cdc
//
// C -> ded
//
//  1. 一个句子应该至少有一个结构体或接口，该结构体标识一种语法、可导出。可以定义其他结构体以辅助描述该结构体，比如 udf 的方法
//  2. 对于非终结符 B 可推导出多个句子，可知 B 的定义对应接口，
//     2.1 这样才可以使得 C 或者 cdc 可以通过内嵌字段 和 实现 B 接口的方法的形式来继承B
//     2.2 继承B的原因是可以将 C 或者 cdc 对应的结构体归纳为 B
//  3. C 作为 B -> C 中的一个句子，因为其可推导出的句子只有 ded，所以其对应的定义可以是一个结构体。
//
// 构造出得ast在vm中进行执行是，应该只是可读取、但是不可修改的，所以这里设置了一些getter方法，对于集合的导出也是保护性拷贝
package ast

// Node 语法的开始字符
type Node interface {
	node()
}

type Expression interface {
	Node
	expression()
}

type UnaryExpression struct {
	op  OperatorNode
	exp Expression
}

func (unary *UnaryExpression) Op() OperatorNode {
	return unary.op
}

func (unary *UnaryExpression) Exp() Expression {
	return unary.exp
}

func (*UnaryExpression) node()       {}
func (*UnaryExpression) expression() {}

// BinaryExpression 关于二元表达式结构体定义的思考
//  1. 如果数组为空，则在编译的时候应该将其“上升”为 变量、常量等
//  2. arguments 中的元算符的优先级和结核性应该是相同的，
//     这样才能保证在计算的时候可以通过 从头到尾遍历的方式(左结合) 或者 从尾到头的方式(右结合) 进行计算
type BinaryExpression struct {
	left      Expression
	priority OperatorPriority
	// todo 结合性
	arguments []binaryExpArgument
}

func (*BinaryExpression) node()       {}
func (*BinaryExpression) expression() {}

func (funcExp *BinaryExpression) Left() Expression {
	return funcExp.left
}

// GetArguments 返回保护性拷贝
func (funcExp *BinaryExpression) GetArguments() []binaryExpArgument {
	forCopy := make([]binaryExpArgument, len(funcExp.arguments))
	copy(forCopy, funcExp.arguments)
	return forCopy
}

type binaryExpArgument struct {
	op  OperatorNode
	arg Expression
}

func (arg *binaryExpArgument) GetOperator() OperatorNode {
	return arg.op
}

func (arg *binaryExpArgument) GetArg() Expression {
	return arg.arg
}

type FuncExpression struct {
	funcName  funcNameNode
	lParen    ControlNode
	arguments []Expression
	rParen    ControlNode
}

func (*FuncExpression) node()                       {}
func (*FuncExpression) expression()                 {}
func (*FuncExpression) atomic()                     {}
func (funcExp *FuncExpression) GetFuncName() string { return funcExp.funcName.name }
func (funcExp *FuncExpression) GetArguments() []Expression {
	forCopy := make([]Expression, len(funcExp.arguments))
	copy(forCopy, funcExp.arguments)
	return forCopy
}

// FuncNameNode 每个元素都搞个node的好处是方便管理扩展，比如添加位置、注释信息等
type funcNameNode struct {
	name string
}
func (*funcNameNode) node()       {}
func (*funcNameNode) expression() {}


type EmptyExpression struct {
}

func (*EmptyExpression) node()       {}
func (*EmptyExpression) expression() {}

type Atomic interface {
	Expression
	atomic()
}

type VariableNode struct {
	name string
}

func (*VariableNode) node()                    {}
func (*VariableNode) atomic()                  {}
func (*VariableNode) expression()              {}
func (variable *VariableNode) GetName() string { return variable.name }

type StringNode struct {
	value string
}

func (*StringNode) node()       {}
func (*StringNode) atomic()     {}
func (*StringNode) expression() {}

func (node *StringNode) GetStringValue() string {
	if node == nil {
		return ""
	}
	return node.value[1 : len(node.value)-1]
}

type NumberNode struct {
	Value int64
}

func (*NumberNode) node()       {}
func (*NumberNode) atomic()     {}
func (*NumberNode) expression() {}

// OperatorNode
// 运算符
type OperatorNode struct {
	op string
	// 这个字段在表达式引擎计算中是没有用的
	// 因为 ast 的结构可以标识运算的优先级
	// 当前赋值主要用于debug和验证程序的正确性，后期应该删除
	priority OperatorPriority

}
func (*OperatorNode) node()   {}
func (*OperatorNode) atomic() {}
func (*OperatorNode) expression() {}
func (operatorNode *OperatorNode) GetOperator() string {
	return operatorNode.op
}

// ControlNode
// 括号，"(",")"
type ControlNode struct {
	value string
}

func (*ControlNode) node()   {}
func (*ControlNode) atomic() {}

// SubNode
//Atomic -> variable
//
//	| string
//	| number
//	| Func
//	| '(' Expression ')' note
type SubNode struct {
	lParen  ControlNode
	subNode Expression
	rParen  ControlNode
}

func (*SubNode) node()       {}
func (*SubNode) atomic()     {}
func (*SubNode) expression() {}
func (n *SubNode) SubNode() Expression {
	return n.subNode
}