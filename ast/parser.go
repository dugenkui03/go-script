package ast

import (
	"errors"
	"fmt"
	"strconv"
)

func Parse(exp string) (Expression, error) {
	tokens, err := getAllTokens(exp)
	if err != nil {
		return nil, err
	}

	// 在当前的语法分析中，空白字符token没有任何作用，所以将空白token移除
	tokensWithoutWhiteToken := filterWhiteToken(tokens)
	if len(tokensWithoutWhiteToken) == 0 {
		return &EmptyExpression{}, nil
	}

	parserPtr := newParser(tokensWithoutWhiteToken)
	return parserPtr.parseInternal()
}

func newParser(tokens []Token) *parser {
	return &parser{
		tokens:  tokens,
		scanner: Scanner{source: tokens},
	}
}

func filterWhiteToken(tokens []Token) []Token {
	tokensWithoutWhiteToken := make([]Token, 0)

	for _, token := range tokens {
		if token.kind != WhiteSpace {
			tokensWithoutWhiteToken = append(tokensWithoutWhiteToken, token)
		}
	}
	return tokensWithoutWhiteToken
}

type parser struct {
	tokens  []Token
	scanner Scanner
}

func (p *parser) parseInternal() (Expression, error) {
	expression, err := p.parseExpression()
	if err != nil {
		return nil, err
	}

	// note 判断tokens是否全部遍历完，
	//		对 有效表达式+无效表达式 的情况做判断
	//		比如 1*2abc 中 abc是有效表达式 1*2后多余的部分
	if p.scanner.peek() != nil {
		return nil, errors.New(fmt.Sprintf("expression before 'line:%d column%d' is valid, the rest is redundant",
			p.scanner.peek().line, p.scanner.peek().column))
	}

	return expression, nil
}

// parseExpression  note 返回值用接口就行，不用接口指针
//
// expression
//
//	: binary
//	;
//
// todo 常量折叠： 1+2 -> 3； -3 -> (-3)；折叠的时候也需要计算，比如数字想加或者字符串拼接，所以不适合在 parser 中进行
func (p *parser) parseExpression() (Expression, error) {
	return p.parseBinaryExpression(firstLevelOp)
}

// parseLevel1BinaryExpression
//
// ```
// binary
//	: level2_binary (First_level_op level2_binary)*
//	;
//
// ```
func (p *parser) parseBinaryExpression(priority OperatorPriority) (Expression, error) {
	// highest priority binary op
	//
	// level2_binary
	//	: signedAtom (Second_level_op signedAtom)*
	//	;
	if priority.isHighestLevelOp() {
		return p.parseSignedAtomBinaryExpression(priority)
	}

	// other priority binary op
	//
	//
	left, err := p.parseBinaryExpression(priority.getIncrement())
	if err != nil {
		return nil, err
	}

	next := p.scanner.peek() // not pop
	if next == nil || !operatorByPriority[priority][next.value] {
		return left, nil
	}

	binaryExpression := BinaryExpression{priority: priority}
	binaryExpression.left = left
	for next != nil && operatorByPriority[priority][next.value] {
		expArgument, err := p.parseBinaryExpArgument(priority)
		if err != nil {
			return nil, err
		}

		binaryExpression.arguments = append(binaryExpression.arguments, *expArgument)
		next = p.scanner.peek() // not pop
	}

	return &binaryExpression, nil
}

// parseSignedAtomicBinaryExpression
//
// ```
// level2_binary
//	: signedAtom (Second_level_op signedAtom)*
//	;
//
// ```
func (p *parser) parseSignedAtomBinaryExpression(priority OperatorPriority) (Expression, error) {
	signedAtom, err := p.parseSignedAtom()
	if err != nil {
		return nil, err
	}

	next := p.scanner.peek()
	if next == nil || !operatorByPriority[priority][next.value] {
		return signedAtom, nil
	}

	atomBinaryExpression := BinaryExpression{priority: priority}
	atomBinaryExpression.left = signedAtom

	for next != nil && operatorByPriority[priority][next.value] {
		expArgument, err := p.parseBinaryExpArgument(priority)

		if err != nil {
			return nil, err
		}
		atomBinaryExpression.arguments = append(atomBinaryExpression.arguments, *expArgument)

		next = p.scanner.peek()
	}

	return &atomBinaryExpression, nil
}

// parseSignedAtom
//
// ```
// signedAtom
//
//	: unary
//	| atomic
//	;
//
// ```
func (p *parser) parseSignedAtom() (Expression, error) {
	lookAhead := p.scanner.peek()

	if unaryOperator[lookAhead.value] {
		return p.parseUnaryExpression()
	}

	return p.parseAtom()
}

// parseUnaryExpression
//
// ```
// unary
//
//	:UnaryOp atom
//	;
//
// ```
func (p *parser) parseUnaryExpression() (Expression, error) {

	operator, err := p.parseUnaryOpe()
	if err != nil {
		return nil, err
	}

	expression := UnaryExpression{}
	expression.op = *operator

	atom, err := p.parseAtom()
	if err != nil {
		return nil, err
	}
	expression.exp = atom

	return &expression, nil
}

// parseUnaryOpe
//
// ```
// UnaryOp
//
//	: '+'
//	| '-'
//	;
//
// ```
func (p *parser) parseUnaryOpe() (*OperatorNode, error) {
	opeToken := p.scanner.pop()

	if _, ok := unaryOperator[opeToken.value]; !ok {
		errorMsg := fmt.Sprintf(
			"expected expression started token instead of '%s'. line:%d, column:%d",
			opeToken.value, opeToken.line, opeToken.column,
		)
		return nil, errors.New(errorMsg)
	}

	return &OperatorNode{
		op: opeToken.value,
	}, nil
}

// parseAtom note 走到这里的时候预期不为空、即不反悔emptyNode，因为上层已经校验过了
//
// ```
// atom
//	: Variable
//	| String
//	| Number
//	| func
//	| sub_node
//	;
//
// ```
func (p *parser) parseAtom() (Expression, error) {
	lookAHead := p.scanner.peek()

	if lookAHead.kind == Variable {
		return p.parseVariable()
	}

	if lookAHead.kind == String {
		return p.parseString()
	}

	if lookAHead.kind == Number {
		return p.parseNumber()
	}

	if lookAHead.kind == Func {
		return p.parseFuncExpression()
	}

	// sub_node LParen expression RParen 的前看符号
	if lookAHead.kind == Control && lookAHead.value == "("{
		return p.parseSubNode()
	}

	errorMsg := fmt.Sprintf(
		"expected Atomic token instead of '%s'. line:%d, column:%d",
		lookAHead.value, lookAHead.line, lookAHead.column,
	)
	return nil, errors.New(errorMsg)
}

// parseFuncExpression
//
// ```
// function
//
//	: Func_name LParen node (Comma node)*  RParen
//	;
//
// ```
func (p *parser) parseFuncExpression() (Expression, error) {

	funcNameNode, err := p.parseFuncName()
	if err != nil {
		return nil, err
	}

	expression := FuncExpression{}
	expression.funcName = *funcNameNode

	lParenNode, err := p.parseLParen()
	if err != nil {
		return nil, err
	}
	expression.lParen = *lParenNode

	//  没有参数的函数
	next := p.scanner.peek()
	if next.kind == Control && next.value == ")" {
		rParentNode, err := p.parseRParen()
		if err != nil {
			return nil, err
		}
		expression.rParen = *rParentNode
		return &expression, nil
	} else {
		firstArg, err := p.parseExpression()
		if err != nil {
			return nil, err
		}

		expression.arguments = append(expression.arguments, firstArg)

		next := p.scanner.peek()
		for next.kind == Comma {
			p.scanner.pop() // swallow comma
			node, err := p.parseExpression()
			if err != nil {
				return nil, err
			}
			expression.arguments = append(expression.arguments, node)
			next = p.scanner.peek()
		}

		rParen, err := p.parseRParen()
		if err != nil {
			return nil, err
		}
		expression.rParen = *rParen

		return &expression, nil
	}
}

func (p *parser) parseVariable() (*VariableNode, error) {
	token := p.scanner.pop()

	return &VariableNode{
		name: token.value,
	}, nil
}

// parseString
//
//	应该在计算的时候考虑？
func (p *parser) parseString() (*StringNode, error) {
	token := p.scanner.pop()
	// assert token type is variable

	return &StringNode{
		value: token.value,
	}, nil
}

func (p *parser) parseNumber() (*NumberNode, error) {
	token := p.scanner.pop()
	num, err := strconv.ParseInt(token.value, 10, 64)
	if err != nil {
		return nil, errors.New("invalid number token")
	}
	return &NumberNode{
		Value: num,
	}, nil
}

func (p *parser) parseFuncName() (*funcNameNode, error) {
	funcNameToken := p.scanner.pop()
	if funcNameToken == nil || funcNameToken.kind != Func {
		errorMsg := fmt.Sprintf(
			"expected function name instead of '%s'. line:%d, column:%d",
			funcNameToken.value, funcNameToken.line, funcNameToken.column,
		)
		return nil, errors.New(errorMsg)
	}

	return &funcNameNode{
		name: funcNameToken.value,
	}, nil
}

func (p *parser) parseLParen() (*ControlNode, error) {
	controlToken := p.scanner.pop()
	if controlToken == nil || controlToken.kind != Control || controlToken.value != "(" {
		errorMsg := fmt.Sprintf(
			"expected left paren instead of '%s'. line:%d, column:%d",
			controlToken.value, controlToken.line, controlToken.column,
		)
		return nil, errors.New(errorMsg)
	}

	return &ControlNode{
		value: controlToken.value,
	}, nil
}

func (p *parser) parseRParen() (*ControlNode, error) {
	controlToken := p.scanner.pop()
	if controlToken == nil || controlToken.kind != Control || controlToken.value != ")" {
		errorMsg := fmt.Sprintf(
			"expected right paren instead of '%s'. line:%d, column:%d",
			controlToken.value, controlToken.line, controlToken.column,
		)
		return nil, errors.New(errorMsg)
	}

	return &ControlNode{
		value: controlToken.value,
	}, nil
}

// parseBinaryExpArgument
// ```
// binary
//
//	: level2_binary (First_level_op level2_binary)*
//	;
//
// or
//
// level2_binary
//
//	: signedAtom (Second_level_op signedAtom)*
//	;
//
// ```
func (p *parser) parseBinaryExpArgument(priority OperatorPriority) (*binaryExpArgument, error) {
	op, err := p.parseOperator(priority)
	if err != nil {
		return nil, err
	}

	argument := binaryExpArgument{}
	argument.op = *op

	var arg Expression
	if priority.isHighestLevelOp() {
		arg, err = p.parseSignedAtom()
	} else {
		arg, err = p.parseBinaryExpression(priority.getIncrement())
	}

	if err != nil {
		return nil, err
	}
	argument.arg = arg

	return &argument, nil
}

// parseSubNode
//
// ```
// sub_node
//	: LParen expression RParen
//	;
//
// ```
func (p *parser) parseSubNode() (Expression, error) {
	_, lErr := p.parseLParen()
	if lErr != nil {
		return nil, lErr
	}

	node, nErr := p.parseExpression()
	if nErr != nil {
		return nil, nErr
	}

	_, rErr := p.parseRParen()
	if rErr != nil {
		return nil, rErr
	}

	return node, nil
}


// parseOperator 获取指定优先级的运算符
func (p *parser) parseOperator(priority OperatorPriority) (*OperatorNode, error) {
	opeToken := p.scanner.pop()
	if _, ok := operatorByPriority[priority][opeToken.value]; !ok {
		errorMsg := fmt.Sprintf(
			"expected %d level op token instead of '%s'. line:%d, column:%d",
			priority, opeToken.value, opeToken.line, opeToken.column,
		)
		return nil, errors.New(errorMsg)
	}

	return &OperatorNode{
		op:       opeToken.value,
		priority: priority,
	}, nil
}
