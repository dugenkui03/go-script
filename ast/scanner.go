package ast

type Scanner struct {
	source []Token // 原始数据
	offset int     // 当前扫描字符的偏移量, 从0开始
}

func (s *Scanner) peek() *Token {
	if s.offset == len(s.source) {
		return nil
	}

	return &s.source[s.offset]
}

func (s *Scanner) pop() *Token {
	if s.offset == len(s.source) {
		return nil
	}

	token := s.source[s.offset]
	s.offset = s.offset + 1
	return &token
}
