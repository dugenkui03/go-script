package ast

import "fmt"

type tokenKind int16

const (
	EOF = iota
	Operator
	Func
	Number
	String
	Variable
	Control // 优先级控制， (, )
	NewLine
	WhiteSpace
	Comma // ,
	// todo Bool
	//todo 三元组，[] 数组
)

var tokenKindDesc = map[tokenKind]string{
	EOF:        "EOF",
	Operator:   "Operator",
	Func:       "Func",
	Number:     "Number",
	String:     "String",
	Variable:   "Variable",
	Control:    "Control",
	NewLine:    "NewLine",
	WhiteSpace: "WhiteSpace",
	Comma:      "Comma",
}

func (kind *tokenKind) String() string {
	if kind == nil {
		return "nil"
	}

	if desc, ok := tokenKindDesc[*kind]; ok {
		return desc
	}

	return "invalid type"
}

type Token struct {
	kind   tokenKind
	value  string
	line   int
	column int
}

func (token *Token) String() string{
	return fmt.Sprintf("{ kind: %d, value: %s, line: %d, column: %d }", token.kind, token.value, token.line, token.column)
}

type expectedLevel int

const (
	first = iota + 1
	second
)
