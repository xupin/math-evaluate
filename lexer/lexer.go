package lexer

import (
	"errors"
	"fmt"
	"unicode"

	"github.com/xupin/math-evaluate/enums"
)

type Lexer struct {
	Expression string
	char       byte
	pos        int
}

func (r *Lexer) Lex() ([]*Token, error) {
	tokens := make([]*Token, 0)
	if len(r.Expression) == 0 {
		return tokens, errors.New("the token list is empty")
	}
	if r.IsChinese() {
		return tokens, errors.New("the token list contains Chinese characters")
	}
	r.char = r.Expression[0]
	for r.pos < len(r.Expression) {
		token := r.Scan()
		if token.t == enums.ILLEGAL {
			return []*Token{}, fmt.Errorf("'%s' is not supported", token.str)
		}
		tokens = append(tokens, token)
	}
	return tokens, nil
}

// 是否包含中文
func (r *Lexer) IsChinese() bool {
	for _, c := range r.Expression {
		if (c >= 65281 && c <= 65374) || c == 12288 {
			return true
		}
	}
	return false
}

// 获取下一个token
func (r *Lexer) Scan() *Token {
	var token *Token
	pos := r.pos
	switch r.char {
	case '(':
		token = &Token{
			str:   string(r.char),
			t:     enums.LPAREN,
			start: pos,
			end:   pos,
		}
		r.NextChar()
	case ')':
		token = &Token{
			str:   string(r.char),
			t:     enums.RPAREN,
			start: pos,
			end:   pos,
		}
		r.NextChar()
	case '{':
		token = &Token{
			str:   string(r.char),
			t:     enums.LBRACE,
			start: pos,
			end:   r.pos,
		}
		r.NextChar()
	case '}':
		token = &Token{
			str:   string(r.char),
			t:     enums.RBRACE,
			start: pos,
			end:   r.pos,
		}
		r.NextChar()
	case ',':
		token = &Token{
			str:   string(r.char),
			t:     enums.COMMA,
			start: pos,
			end:   pos,
		}
		r.NextChar()
	case
		'0',
		'1',
		'2',
		'3',
		'4',
		'5',
		'6',
		'7',
		'8',
		'9':
		for r.IsDigit() {
			if !r.NextChar() {
				break
			}
		}
		token = &Token{
			str:   string(r.Expression[pos:r.pos]),
			t:     enums.NUMBER,
			start: pos,
			end:   r.pos - 1,
		}
	case '+':
		token = &Token{
			str:   string(r.char),
			t:     enums.ADD,
			start: pos,
			end:   pos,
		}
		r.NextChar()
	case '-':
		token = &Token{
			str:   string(r.char),
			t:     enums.SUB,
			start: pos,
			end:   pos,
		}
		r.NextChar()
	case '*':
		token = &Token{
			str:   string(r.char),
			t:     enums.MUL,
			start: pos,
			end:   pos,
		}
		if r.NextChar() && r.char == '*' {
			token = &Token{
				str:   "**",
				t:     enums.XOR,
				start: pos,
				end:   pos,
			}
			r.NextChar()
		}
	case '/':
		token = &Token{
			str:   string(r.char),
			t:     enums.QUO,
			start: pos,
			end:   pos,
		}
		r.NextChar()
	case '%':
		token = &Token{
			str:   string(r.char),
			t:     enums.REM,
			start: pos,
			end:   pos,
		}
		r.NextChar()
	case '^':
		token = &Token{
			str:   string(r.char),
			t:     enums.XOR,
			start: pos,
			end:   pos,
		}
		r.NextChar()
	default:
		if r.IsWs() { // 跳过空白字符
			for r.NextChar() {
				if !r.IsWs() {
					break
				}
			}
			return r.Scan()
		} else if r.IsLetter() { // 判断是不是字母（函数）
			for r.IsLetter() {
				if !r.NextChar() {
					break
				}
			}
			token = &Token{
				str:   string(r.Expression[pos:r.pos]),
				t:     enums.FUNC,
				start: pos,
				end:   r.pos - 1,
			}
		} else if r.IsVar() { // 判断是不是变量
			for r.IsVar() {
				if !r.NextChar() {
					break
				}
			}
			token = &Token{
				str:   string(r.Expression[pos:r.pos]),
				t:     enums.VAR,
				start: pos,
				end:   r.pos - 1,
			}
		} else {
			token = &Token{
				str:   string(r.char),
				t:     enums.ILLEGAL,
				start: pos,
				end:   r.pos,
			}
		}
	}
	return token
}

// 下一个字符
func (r *Lexer) NextChar() bool {
	// 判断是否越界
	r.pos++
	if r.pos >= len(r.Expression) {
		return false // eof
	}
	r.char = r.Expression[r.pos] // 移动[当前字符位置]
	return true
}

// 空白字符
func (r *Lexer) IsWs() bool {
	return unicode.IsSpace(rune(r.char))
}

// 数字（小数）
func (r *Lexer) IsDigit() bool {
	return unicode.IsNumber(rune(r.char)) || r.char == '.'
}

// 字母
func (r *Lexer) IsLetter() bool {
	return r.char >= 'a' && r.char <= 'z'
}

// 变量
func (r *Lexer) IsVar() bool {
	return r.char >= 'A' && r.char <= 'Z'
}
