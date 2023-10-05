package goscript

import (
	"github.com/stretchr/testify/assert"
	"goscript/vm"
	"testing"
)

// todo 计算中报错，比如将 string 传入参数类型为 int 的udf
var resultByValidExp = map[string]interface{}{
	"1+2+3-4": 2,
	"-3*2+1 ": -5,
	"-1-1":    -2,
	"1+2*3":   7, // 优先级
	"(1+2)*3": 9, // 括号优先级
	"a":       2,
	"oneNotFold()":1,
	"oneNotFold()+1":2,
	"a+2*3":      8,
	"2+a1*a":      2,
	"1+one()":    2,
	"1+same(a)":  3,
	"1+same(100)":  101,
	"a+b":5,
	"a+b*c+same(a)":16,
	"-same(6)+5":-1,
}

//	"-(a)",
//	"(a)",
//	"-(a+b)",
//	"(1+2)*3",
//	"1+2+3-4",
//	"1+2+3",
//	"1+2*3+4",
//	"1+2*3",
//	"(1+2)*3",
//	"test()",
//	"2 * -3",
//	"-3 * 2 + 1",
//	"test(1+2,a,test())",
//	"123+ab*100*\"abc\"",
//	"123+ab*100*\"ab\\\"c\"",
//	"test(a)",
//	"test(a)+1",
//	"a+test(a)",
//	"1+same(100)",

// - `t.Errorf`：报告测试失败，并输出错误信息。
//- `t.Fatalf`：报告测试失败，并输出错误信息，然后终止测试。
//- `t.Logf`：输出日志信息，不会导致测试失败。
//- `t.Fail`：标记测试为失败，但继续执行测试。
//- `t.FailNow`：标记测试为失败，并立即终止测试。

var virtualMachine = vm.NewVM()

func init() {
	_ = virtualMachine.RegisterFunc0("oneNotFold", true, func() (interface{}, error) {
		return 1, nil
	})

	_ = virtualMachine.RegisterFunc0("one", true, func() (interface{}, error) {
		return 1, nil
	})

	_ = virtualMachine.RegisterFunc1("same", true, func(arg vm.Value) (interface{}, error) {
		intValue, _ := arg.AsInt()
		return intValue, nil
	})
}

func TestEval(t *testing.T) {
	env := map[string]interface{}{"a": 2,"b":3,"c":4}
	for exp, expected := range resultByValidExp {
		eval, err := virtualMachine.Eval(exp, env)
		if err != nil {
			t.Fatalf("ext: '%s', err: %v", exp, err)
			continue
		}
		result, _ := eval.AsInt()
		assert.Equal(t, expected, result, "exp: "+exp)
	}
}



//func TestEvalBadCase(t *testing.T) {
//	exp := "-1-1"
//	expected := int64(-2)
//	eval, err := Eval("-1-1", nil)
//	if err != nil {
//		t.Fatalf("ext: '%s', err: %v", exp, err)
//		return
//	}
//
//	assert.Equal(t, expected, eval, "exp: "+exp)
//}

