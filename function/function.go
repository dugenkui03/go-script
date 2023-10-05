package function

func NewFunction(name string, f interface{}, argumentsNum int, false bool) Function {
	return Function{
		name:         name,
		f:            f,
		argumentsNum: argumentsNum,
		allowFold:    false,
	}
}

type Function struct {
	name string

	// func
	f interface{}

	// number of arguments
	argumentsNum int

	// 如果表达式中函数参数是常量，是否允许对结果进行预计算并替换表达式中的函数调用部分
	// eg: a+add(1,2) -> a+3
	allowFold bool
}

func (f Function) Name() string {
	return f.name
}

func (f Function) AllowFold() bool {
	return f.allowFold
}

func (f Function) F() interface{} {
	return f.f
}

func (f Function) ArgumentsNum() int {
	return f.argumentsNum
}
