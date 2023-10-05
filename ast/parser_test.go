package ast

import (
	"github.com/stretchr/testify/assert"
	"testing"
)



// - `t.Errorf`：报告测试失败，并输出错误信息。
//- `t.Fatalf`：报告测试失败，并输出错误信息，然后终止测试。
//- `t.Logf`：输出日志信息，不会导致测试失败。
//- `t.Fail`：标记测试为失败，但继续执行测试。
//- `t.FailNow`：标记测试为失败，并立即终止测试。

var validExpressions = []string{
	"a+b",
	"a+b*c*d+e+test(f)",
	"-(a)",
	"(a)",
	"-(a+b)",
	"(1+2)*3",
	"1+2+3-4",
	"1+2+3",
	"1+2*3+4",
	"1+2*3",
	"(1+2)*3",
	"test()",
	"2 * -3",
	"-3 * 2 + 1",
	"test(1+2,a,test())",
	"123+ab*100*\"abc\"",
	"123+ab*100*\"ab\\\"c\"",
	"test(a)",
	"test(a)+1",
	"a+test(a)",
	"1+same(100)",
}

var invalidExpressions = []string{
	`123ab+cd`,
	`123ab+[cd`,
	"123ab",
	"1+",
	"test(a,)",
}


func TestParse(t *testing.T) {
	for _,exp := range validExpressions {
		expression, err := Parse(exp)
		assert.Nil(t, err)
		assert.NotNil(t, expression)
		t.Logf("\nsource code: "+exp)
		PrintAST(expression)
	}
}

func TestParseSingleBadCase(t *testing.T){
	expression, err := Parse("1+same(100)")
	assert.Nil(t, err)
	assert.NotNil(t, expression)
}

func TestParseBadCase(t *testing.T) {
	for _, exp := range invalidExpressions {
		expression, err := Parse(exp)
		t.Logf("%v",err)
		assert.NotNil(t, err)
		assert.Nil(t, expression)
	}
}
