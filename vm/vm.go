package vm

import (
	"errors"
	"fmt"
	"goscript/ast"
	"goscript/function"
)

var defaultVM = &VM{
	expressionCache: make(map[string]ast.Expression),
	funcByName:      make(map[string]function.Function),
}

// Eval 用户大多数情况使用的还是默认的 vm
func Eval(exp string, env map[string]interface{}) (*Value, error) {
	return defaultVM.Eval(exp, env)
}

func NewVM() *VM {
	return &VM{
		expressionCache: make(map[string]ast.Expression),
		funcByName:      make(map[string]function.Function),
	}
}

type VM struct {
	expressionCache map[string]ast.Expression
	funcByName      map[string]function.Function
}

func (vm *VM) Eval(exp string, env map[string]interface{}) (*Value, error) {
	cachedExp := vm.getExpressionFromCache(exp)
	if cachedExp != nil {
		return vm.calInternal(cachedExp, env)
	}

	expression, err := ast.Parse(exp)
	if err != nil {
		return nil, err
	}

	expression, err = ast.Optimize(expression, vm.funcByName)
	if err != nil {
		return nil, fmt.Errorf("occur error when optimize expression:%v", err)
	}
	vm.setExpressionCache(exp, expression)

	return vm.calInternal(expression, env)
}

func (vm *VM) calInternal(exp ast.Expression, env map[string]interface{}) (*Value, error) {
	rawValue, err := vm.cal(exp, env)
	if err != nil {
		return nil, err
	}

	return &Value{rawValue: rawValue}, nil
}

func (vm *VM) cal(exp ast.Expression, env map[string]interface{}) (interface{}, error) {
	switch expression := exp.(type) {
	case *ast.EmptyExpression:
		return nil, nil
	case *ast.NumberNode:
		return expression.Value, nil
	case *ast.StringNode:
		return expression.GetStringValue(), nil
	case *ast.VariableNode:
		return vm.calVariable(expression.GetName(), env)
	case *ast.BinaryExpression:
		return vm.calBinary(*expression, env)
	case *ast.UnaryExpression:
		return vm.calUnary(*expression, env)
	case *ast.FuncExpression:
		return vm.calFuncExpression(*expression, env)
	case *ast.SubNode:
		return vm.cal(expression.SubNode(), env)
	default:
		return nil, errors.New(fmt.Sprintf("invalid expression type %T", expression))
	}
}

func (vm *VM) calFuncExpression(expression ast.FuncExpression, env map[string]interface{}) (interface{}, error) {
	// note 编译的时候就应该判断一下有没有udf

	var f function.Function
	if val, ok := vm.funcByName[expression.GetFuncName()]; !ok {
		return nil, errors.New(fmt.Sprintf("invalid udf named '%s'", expression.GetFuncName()))
	} else {
		f = val
	}

	// note 应该编译的时候就发现这个问题，至少应该打印个warn日志（万一用户只想编译做分析、不想执行？）
	// note 或者编译的时候可以让用户可选的忽略udf的是否存在的检查？
	if f.ArgumentsNum() != -1 && f.ArgumentsNum() != len(expression.GetArguments()) {
		return nil, errors.New(fmt.Sprintf("the func of '%s' require %d argument instead of %d",
			f.Name(), f.ArgumentsNum(), len(expression.GetArguments())))
	}

	args := make([]Value, 0)
	for _, argNode := range expression.GetArguments() {
		rawValue, err := vm.cal(argNode, env)
		if err != nil {
			return nil, err
		}
		args = append(args, Value{rawValue})
	}

	return calUdf(f, args)
}

func calUdf(f function.Function, args []Value) (interface{}, error) {
	// todo 有可能panic，在哪捕获panic比较好？
	switch f.ArgumentsNum() {
	case 0:
		if ff, ok := f.F().(func() (interface{}, error)); ok {
			return ff()
		}
	case 1:
		if ff, ok := f.F().(func(arg Value) (interface{}, error)); ok {
			return ff(args[0])
		}
	case 2:
		if ff, ok := f.F().(func(arg1, arg2 Value) (interface{}, error)); ok {
			return (ff)(args[0], args[1])
		}

	case 3:
		if ff, ok := f.F().(func(arg1, arg2, arg3 Value) (interface{}, error)); ok {
			return (ff)(args[0], args[1], args[2])
		}
	default:
		return nil, errors.New("todo: 使用反射或者生成代码")
	}

	return nil, errors.New("should not invoke here")
}

func (vm *VM) calBinary(exp ast.BinaryExpression, env map[string]interface{}) (interface{}, error) {
	firstVal, err := vm.cal(exp.Left(), env)
	if err != nil {
		return nil, err
	}

	tmpResult := firstVal
	for _, argument := range exp.GetArguments() {
		argumentVal, aErr := vm.cal(argument.GetArg(), env)
		if aErr != nil {
			return nil, aErr
		}

		// note 左结合的运算
		oResult, oerr := vm.opeCal(argument.GetOperator(), tmpResult, argumentVal)
		if oerr != nil {
			return nil, oerr
		}
		tmpResult = oResult
	}

	return tmpResult, nil
}

func (vm *VM) calUnary(unaryExpression ast.UnaryExpression, env map[string]interface{}) (interface{}, error) {
	expValue, err := vm.cal(unaryExpression.Exp(), env)
	if err != nil {
		return nil, err
	}
	numberValue, err := Value{rawValue: expValue}.AsInt64()
	if err != nil {
		return nil, err
	}

	op := unaryExpression.Op()
	operator := (&op).GetOperator()

	if operator == "-" {
		return -numberValue, nil
	} else if operator == "+" {
		return numberValue, nil
	} else {
		return nil, errors.New("invalid unary operator '" + operator + "'")
	}
}

// todo
//  1. 返回结果的包装类，可能包括结果类型、值以及获取转换后类型值的方法等
//  2. 变量替换成参数
func (vm *VM) opeCal(op ast.OperatorNode, arg1, arg2 interface{}) (interface{}, error) {
	switch op.GetOperator() {
	// todo 操作符和具体函数的绑定关系
	case "+":
		i, err := function.Int64(arg1)
		if err != nil {
			return nil, err
		}

		i2, err := function.Int64(arg2)
		if err != nil {
			return nil, err
		}

		return interface{}(function.Add(i, i2)), nil
	case "-":
		i, err := function.Int64(arg1)
		if err != nil {
			return nil, err
		}

		i2, err := function.Int64(arg2)
		if err != nil {
			return nil, err
		}

		return interface{}(i - i2), nil
	case "*":
		i, err := function.Int64(arg1)
		if err != nil {
			return nil, err
		}

		i2, err := function.Int64(arg2)
		if err != nil {
			return nil, err
		}

		return interface{}(i * i2), nil

	case "/":
		i, err := function.Int64(arg1)
		if err != nil {
			return nil, err
		}

		i2, err := function.Int64(arg2)
		if err != nil {
			return nil, err
		}

		return interface{}(i / i2), nil

	case "%":
		i, err := function.Int64(arg1)
		if err != nil {
			return nil, err
		}

		i2, err := function.Int64(arg2)
		if err != nil {
			return nil, err
		}

		return interface{}(i % i2), nil
	default:
		return nil, errors.New("invalid operator:" + op.GetOperator())
	}
}

func (vm *VM) calVariable(variableName string, env map[string]interface{}) (interface{}, error) {
	// note 如果表达式只有一个变量 a，则直接返回a对应的对象，int/int32等也不会返回对应的转换后的值
	// todo
	//		1. 将各种类型的 int 统一为 int64
	//		2. env 不仅可以是map，而且可以是对象、从对象中反射取值。internal.Fetch(env(interface{}),key)
	return env[variableName], nil
}

func (vm *VM) getExpressionFromCache(exp string) ast.Expression {
	if vm == nil {
		return nil
	}

	return vm.expressionCache[exp]
}

func (vm *VM) setExpressionCache(exp string, expression ast.Expression) {
	if vm == nil {
		return
	}

	vm.expressionCache[exp] = expression
}
