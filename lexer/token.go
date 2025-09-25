package lexer

import "github.com/xupin/math-evaluate/enums"

type Token struct {
	str   string
	t     int
	start int
	end   int
}

func EOF() *Token {
	return &Token{
		t: enums.EOF,
	}
}

func (t *Token) GetStr() string {
	return t.str
}

func (t *Token) GetType() int {
	return t.t
}
