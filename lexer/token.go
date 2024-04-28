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
		str: "eof",
		t:   enums.EOF,
	}
}

func (r *Token) GetStr() string {
	return r.str
}

func (r *Token) GetType() int {
	return r.t
}
