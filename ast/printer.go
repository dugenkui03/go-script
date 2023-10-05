package ast

import "fmt"


func PrintAST(exp Expression) {
	if exp == nil {
		fmt.Println("<nil>")
	}

	var printVisitor func(deep int, exp Expression) WalkControl

	printVisitor = func(deep int, exp Expression) WalkControl {
		switch e := exp.(type) {
		case *NumberNode:
			printDeep(deep)
			println(e.Value)
		case *StringNode:
			printDeep(deep)
			println(e.GetStringValue())
		case *VariableNode:
			printDeep(deep)
			println(e.name)
		case *OperatorNode:
			printDeep(deep)
			println(e.op)
		case *funcNameNode:
			printDeep(deep)
			println("funcName: " + e.name)
		case *EmptyExpression:
			printDeep(deep)
			println("<EmptyExpression>")
		case *BinaryExpression:
			printDeep(deep)
			println("<BinaryExpression>")
			//printVisitor(deep+1, e.left)
			//for i, arg := range e.GetArguments() {
			//	printDeep(deep)
			//	println("op: " + arg.GetOperator().GetOperator())
			//	printDeep(deep)
			//	print(fmt.Sprintf("arg %d: \n", i))
			//	printVisitor(deep+1, arg.GetArg())
			//}
		case *UnaryExpression:
			printDeep(deep)
			println("<UnaryExpression>")
			//printDeep(deep)
			//println(fmt.Sprintf("op: %s", e.op.GetOperator()))
			//printVisitor(deep+1, e.exp)
		case *FuncExpression:
			printDeep(deep)
			println("<FuncExpression>")
			//printDeep(deep)
			//println("funcName: %s",e.GetFuncName())
			//for i, arg := range e.GetArguments() {
			//	printDeep(deep)
			//	print(fmt.Sprintf("arg %d: ", i))
			//	printVisitor(deep+1, arg)
			//}
		case *SubNode:
			printDeep(deep)
			println(fmt.Sprintf("<SubNode>"))
			//printVisitor(deep+1, e.subNode)
		}
		return Continue
	}

	WalkDeepFirst(exp, printVisitor)
}

func printDeep(deep int) {
	for i := 0; i < deep-1; i++ {
		print("│\t")
	}
	print("├─")
}