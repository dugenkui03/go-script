package ast

import "fmt"

type WalkControl int

const (
	Continue = iota + 1
	Abort
)

func WalkDeepFirst(exp Expression, f func(deep int, exp Expression) WalkControl) {
	walk(exp, 1, f)
}

func walk(exp Expression, deep int, f func(deep int, exp Expression) WalkControl) {
	if f(deep, exp) == Abort {
		return
	}

	switch e := exp.(type) {
	case *BinaryExpression:
		walk(e.left, deep+1, f)
		for _, arg := range e.GetArguments() {
			walk(&arg.op, deep+1, f)
			walk(arg.GetArg(), deep+1, f)
		}

	case *UnaryExpression:
		walk(&e.op, deep+1, f)
		walk(e.exp, deep+1, f)

	case *FuncExpression:
		walk(&e.funcName, deep+1, f)
		for _, arg := range e.GetArguments() {
			walk(arg, deep+1, f)
		}

	case *SubNode:
		walk(e.subNode, deep+1, f)

	case *EmptyExpression, *NumberNode, *StringNode, *VariableNode, *OperatorNode, *funcNameNode:
	default:
		panic(fmt.Sprintf("ast.Walk: unexpected expression type %T", e))
	}

}
