package ast

import (
	"errors"
	"fmt"
	"unicode"
)

func getAllTokens(exp string) ([]Token, error) {
	lexerPtr := newLexer(exp)
	return lexerPtr.getTokens()
}

func newLexer(exp string) lexer {
	return lexer{
		source: []rune(exp),
		offset: 0,
		lines:  []int{0},
	}
}

type lexer struct {
	source []rune
	lines  []int // 包含每一行的offset
	offset int
}

func (lexer *lexer) getTokens() ([]Token, error) {

	tokens := make([]Token, 0)

	token, err := lexer.getNextToken()
	if err != nil {
		return nil, err
	}
	for token != nil {
		tokens = append(tokens, *token)
		token, err = lexer.getNextToken()
		if err != nil {
			return nil, err
		}
	}

	return tokens, nil
}

// Operator + - * / % 函数
// Number 数字
// Variable 变量
func (lexer *lexer) getNextToken() (*Token, error) {

	chPtr := lexer.getNextRune()
	if chPtr == nil {
		return nil, nil
	}

	ch := *chPtr

	switch {
	case isBasicOperator(ch):
		pos := lexer.offset
		return &Token{
			kind:   Operator,
			value:  string(ch),
			line:   len(lexer.lines),
			column: pos - int(lexer.lines[len(lexer.lines)-1]),
		}, nil
	case unicode.IsDigit(ch):
		// todo 小数、科学计数法
		pos := lexer.offset

		start, end := lexer.offset-1, lexer.offset
		next := lexer.getNextRune()
		for next != nil && unicode.IsDigit(*next) {
			next = lexer.getNextRune()
			end++
		}

		if next != nil {
			lexer.rollbackRune()
		}

		return &Token{
			kind:  Number,
			value: string(lexer.source[start:end]),
			line:  len(lexer.lines),
			// note column 应该从开始字符开始算
			column: pos - lexer.lines[len(lexer.lines)-1],
		}, nil
	case isQuote(ch):
		isSingleQuote := ch == '\''
		pos := lexer.offset
		start, end := lexer.offset-1, lexer.offset
		if isSingleQuote {
			var escape bool = false
			next := lexer.getNextRune()
			for next != nil && (*next != '\'' || escape) {
				// note 用于判断转义
				if *next == '\\' && !escape {
					escape = true
				} else if *next == '\'' {
					escape = false
				}
				next = lexer.getNextRune()
				end++
			}

			if *next == '\'' {
				end++
			}
		} else {
			var escape bool = false
			next := lexer.getNextRune()
			for next != nil && (*next != '"' || escape) {
				if *next == '\\' && !escape {
					escape = true
				} else if *next == '"' {
					escape = false
				}
				next = lexer.getNextRune()
				end++
			}

			if *next == '"' {
				end++
			}
		}
		// todo 结束的时候应该判断是否是正常结束，比如是否有 " 或者 '
		// 最后一个不属于字符串，所以回滚
		if !lexer.scanToEnd() {
			lexer.rollbackRune()
		}
		return &Token{
			kind:   String,
			value:  string(lexer.source[start:end]),
			line:   len(lexer.lines),
			column: pos - int(lexer.lines[len(lexer.lines)-1]),
		}, nil

	case unicode.IsLetter(ch) || ch == '_':
		pos := lexer.offset
		start, end := lexer.offset-1, lexer.offset
		next := lexer.getNextRune()
		for next != nil && (unicode.IsNumber(*next) || unicode.IsLetter(*next) || *next == '_') {
			next = lexer.getNextRune()
			end++
		}

		// todo 如果只有一个 '_' 则判定为非法

		// 扫描到了不符合条件的字符
		if next != nil {
			lexer.rollbackRune()
		}

		// 如果是以字母结尾，或者 即使不是结尾、但字母token后边跟着的不是 (，则该字母所在的字符串是变量
		if next == nil || *next != '(' {
			return &Token{
				kind:   Variable,
				value:  string(lexer.source[start:end]),
				line:   len(lexer.lines),
				column: pos - lexer.lines[len(lexer.lines)-1],
			}, nil
		} else {
			// 如果遇到了括号，则是函数
			return &Token{
				kind:   Func,
				value:  string(lexer.source[start:end]),
				line:   len(lexer.lines),
				column: pos - int(lexer.lines[len(lexer.lines)-1]),
			}, nil
		}

	case isParen(ch):
		pos := lexer.offset
		return &Token{
			kind: Control, value: string(ch),
			line:   len(lexer.lines),
			column: pos - int(lexer.lines[len(lexer.lines)-1]),
		}, nil
	case ch == '\n':
		pos := lexer.offset
		tokenPtr := &Token{
			kind: NewLine, value: string(ch),
			line:   len(lexer.lines),
			column: pos - int(lexer.lines[len(lexer.lines)-1]),
		}
		lexer.lines = append(lexer.lines, lexer.offset)
		return tokenPtr, nil

	case unicode.Is(unicode.White_Space, ch):
		//- 空格字符（U+0020）
		//- 制表符（U+0009）todo ?
		//- 换行符（U+000A）todo ?
		//- 回车符（U+000D）
		//- 垂直制表符（U+000B）
		//- 换页符（U+000C）
		return &Token{
			kind:   WhiteSpace,
			value:  string(ch),
			line:   len(lexer.lines),
			column: lexer.offset - lexer.lines[len(lexer.lines)-1],
		}, nil
	case ch == ',':
		return &Token{
			kind:   Comma,
			value:  string(ch),
			line:   len(lexer.lines),
			column: lexer.offset - lexer.lines[len(lexer.lines)-1],
		}, nil
	default:
		errorMsg := fmt.Sprintf("invalid token '%s':%d:%d:\n %s...",
			string(ch), len(lexer.lines), lexer.offset-lexer.lines[len(lexer.lines)-1], string(lexer.source[0:lexer.offset]))
		return nil, errors.New(errorMsg)
	}
}

func isParen(ch rune) bool {
	return ch == '(' || ch == ')'
}

func (lexer *lexer) getNextRune() *rune {
	// 边界条件，已经遍历完了所有数据
	if lexer.offset == len(lexer.source) {
		return nil
	}

	ch := lexer.source[lexer.offset]
	lexer.offset = lexer.offset + 1
	return &ch
}

func (lexer *lexer) rollbackRune() {
	lexer.offset = lexer.offset - 1
}

func (lexer *lexer) scanToEnd() bool {
	return lexer.offset == len(lexer.source)
}

func isQuote(ch rune) bool {
	return ch == '\'' || ch == '"'
}

func isBasicOperator(ch rune) bool {
	return ch == '+' || ch == '-' || ch == '*' || ch == '/' || ch == '%'
}
