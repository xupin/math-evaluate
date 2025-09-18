package lexer

import (
	"errors"
	"fmt"
	"unicode"

	"github.com/xupin/math-evaluate/enums"
)

type lexer struct {
	expression string
	char       byte
	pos        int
}

func New(exp string) *lexer {
	return &lexer{
		expression: exp,
	}
}

func (l *lexer) Lex() ([]Token, error) {
	tokens := make([]Token, 0)
	if len(l.expression) == 0 {
		return tokens, errors.New("input expression is empty")
	}
	if l.isChinese() {
		return tokens, errors.New("input expression contains Chinese characters")
	}
	l.char = l.expression[0]
	for l.pos < len(l.expression) {
		token := l.scan()
		if token.t == enums.ILLEGAL {
			return []Token{}, fmt.Errorf("'%s' is not supported", token.str)
		}
		tokens = append(tokens, token)
	}
	return tokens, nil
}

// 是否包含中文
func (l *lexer) isChinese() bool {
	for _, c := range l.expression {
		if (c >= 65281 && c <= 65374) || c == 12288 {
			return true
		}
	}
	return false
}

// 获取下一个token
func (l *lexer) scan() Token {
	var token Token
	pos := l.pos
	switch l.char {
	case '(':
		token = Token{
			str:   string(l.char),
			t:     enums.LPAREN,
			start: pos,
			end:   pos,
		}
		l.nextChar()
	case ')':
		token = Token{
			str:   string(l.char),
			t:     enums.RPAREN,
			start: pos,
			end:   pos,
		}
		l.nextChar()
	case '{':
		token = Token{
			str:   string(l.char),
			t:     enums.LBRACE,
			start: pos,
			end:   l.pos,
		}
		l.nextChar()
	case '}':
		token = Token{
			str:   string(l.char),
			t:     enums.RBRACE,
			start: pos,
			end:   l.pos,
		}
		l.nextChar()
	case ',':
		token = Token{
			str:   string(l.char),
			t:     enums.COMMA,
			start: pos,
			end:   pos,
		}
		l.nextChar()
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
		for l.isDigit() {
			if !l.nextChar() {
				break
			}
		}
		token = Token{
			str:   string(l.expression[pos:l.pos]),
			t:     enums.NUMBER,
			start: pos,
			end:   l.pos - 1,
		}
	case '+':
		token = Token{
			str:   string(l.char),
			t:     enums.ADD,
			start: pos,
			end:   pos,
		}
		l.nextChar()
	case '-':
		token = Token{
			str:   string(l.char),
			t:     enums.SUB,
			start: pos,
			end:   pos,
		}
		l.nextChar()
	case '*':
		token = Token{
			str:   string(l.char),
			t:     enums.MUL,
			start: pos,
			end:   pos,
		}
		if l.nextChar() && l.char == '*' {
			token = Token{
				str:   "**",
				t:     enums.XOR,
				start: pos,
				end:   pos,
			}
			l.nextChar()
		}
	case '/':
		token = Token{
			str:   string(l.char),
			t:     enums.QUO,
			start: pos,
			end:   pos,
		}
		l.nextChar()
	case '%':
		token = Token{
			str:   string(l.char),
			t:     enums.REM,
			start: pos,
			end:   pos,
		}
		l.nextChar()
	case '^':
		token = Token{
			str:   string(l.char),
			t:     enums.XOR,
			start: pos,
			end:   pos,
		}
		l.nextChar()
	default:
		if l.isWs() { // 跳过空白字符
			for l.nextChar() {
				if !l.isWs() {
					break
				}
			}
			return l.scan()
		} else if l.isLetter() { // 判断是不是字母（函数）
			for l.isLetter() {
				if !l.nextChar() {
					break
				}
			}
			token = Token{
				str:   string(l.expression[pos:l.pos]),
				t:     enums.FUNC,
				start: pos,
				end:   l.pos - 1,
			}
		} else if l.isVar() { // 判断是不是变量
			for l.isVar() {
				if !l.nextChar() {
					break
				}
			}
			token = Token{
				str:   string(l.expression[pos:l.pos]),
				t:     enums.VAR,
				start: pos,
				end:   l.pos - 1,
			}
		} else {
			token = Token{
				str:   string(l.char),
				t:     enums.ILLEGAL,
				start: pos,
				end:   l.pos,
			}
		}
	}
	return token
}

// 下一个字符
func (l *lexer) nextChar() bool {
	// 判断是否越界
	l.pos++
	if l.pos >= len(l.expression) {
		return false // eof
	}
	l.char = l.expression[l.pos] // 移动[当前字符位置]
	return true
}

// 空白字符
func (l *lexer) isWs() bool {
	return unicode.IsSpace(rune(l.char))
}

// 数字（小数）
func (l *lexer) isDigit() bool {
	return unicode.IsNumber(rune(l.char)) || l.char == '.'
}

// 字母
func (l *lexer) isLetter() bool {
	return l.char >= 'a' && l.char <= 'z'
}

// 变量
func (l *lexer) isVar() bool {
	return l.char >= 'A' && l.char <= 'Z'
}
