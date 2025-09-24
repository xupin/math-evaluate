package parser

import (
	"errors"
	"fmt"
	"log"
	"math"
	"math/rand/v2"
	"strconv"

	"github.com/xupin/math-evaluate/enums"
	"github.com/xupin/math-evaluate/interfaces"
	"github.com/xupin/math-evaluate/lexer"
)

type funcType func(...interfaces.INode) float64

type parser struct {
	tokens   []lexer.Token
	curToken *lexer.Token
	index    int
	err      error
	vars     map[string]float64
	fn       map[string]funcType
}

func New(tokens []lexer.Token) *parser {
	p := &parser{
		tokens: tokens,
		vars:   make(map[string]float64),
		fn:     make(map[string]funcType),
	}
	p.loadFunc()
	if len(tokens) > 0 {
		p.curToken = &tokens[0]
	}
	return p
}

// 内置函数
func (p *parser) loadFunc() {
	p.fn = map[string]funcType{
		"min": func(n ...interfaces.INode) float64 {
			if len(n) < 2 {
				return 0
			}
			return math.Min(n[0].Evaluate(), n[1].Evaluate())
		},
		"max": func(n ...interfaces.INode) float64 {
			if len(n) < 2 {
				return 0
			}
			return math.Max(n[0].Evaluate(), n[1].Evaluate())
		},
		"floor": func(n ...interfaces.INode) float64 {
			if len(n) < 1 {
				return 0
			}
			return math.Floor(n[0].Evaluate())
		},
		"round": func(n ...interfaces.INode) float64 {
			if len(n) < 1 {
				return 0
			}
			return math.Round(n[0].Evaluate())
		},
		"random": func(n ...interfaces.INode) float64 {
			v := rand.Float64()
			log.Printf("random %+v", v)
			return v
		},
	}
}

// 设置变量
func (p *parser) SetVar(key string, value float64) {
	p.vars[key] = value
}

// 设置函数
func (p *parser) SetFunc(name string, fn funcType) {
	p.fn[name] = fn
}

// 解析入口
func (p *parser) Parse() (interfaces.INode, error) {
	if len(p.tokens) == 0 {
		return nil, errors.New("the token list is empty")
	}
	return p.compile(), p.err
}

// 构建表达式树
func (p *parser) compile() interfaces.INode {
	left := p.parseExpr()
	return p.parseRight(1, left)
}

// 解析表达式左侧
func (p *parser) parseExpr() interfaces.INode {
	switch p.curToken.GetType() {
	case enums.NUMBER:
		return p.parseNumber()
	case enums.SUB: // e.g. -1
		p.nextToken()
		return &stmt{
			Type:  enums.SUB,
			Left:  &number{Val: 0},
			Right: p.parseExpr(),
		}
	case enums.LPAREN:
		return p.parseStmt()
	case enums.IDENT:
		return p.parseIdent()
	default:
		p.err = fmt.Errorf("unexpected token '%s'", p.curToken.GetStr())
		return nil
	}
}

// 解析二元运算右侧
func (p *parser) parseRight(priority int, left interfaces.INode) interfaces.INode {
	for {
		curPriority := p.priority()
		if curPriority < priority {
			return left
		}
		op := p.curToken.GetType()
		p.nextToken()
		right := p.parseExpr()
		if right == nil {
			return left
		}
		if curPriority < p.priority() {
			right = p.parseRight(curPriority, right)
			if right == nil {
				return left
			}
		}
		left = &stmt{
			Type:  op,
			Left:  left,
			Right: right,
		}
	}
}

// 解析括号表达式
func (p *parser) parseStmt() interfaces.INode {
	p.nextToken()
	if p.curToken.GetType() == enums.RPAREN {
		p.err = errors.New("empty parentheses expression")
		return nil
	}
	node := p.compile()
	if p.curToken.GetType() != enums.RPAREN {
		p.err = fmt.Errorf("expression expects ')' to close, got '%s'", p.curToken.GetStr())
		return nil
	}
	p.nextToken()
	return node
}

// 解析数字
func (p *parser) parseNumber() interfaces.INode {
	f, _ := strconv.ParseFloat(p.curToken.GetStr(), 64)
	node := &number{Val: f}
	p.nextToken()
	return node
}

// 解析函数、变量
func (p *parser) parseIdent() interfaces.INode {
	ident := p.curToken.GetStr()
	if next := p.peekNextToken(); next != nil && next.GetType() == enums.LPAREN {
		return p.parseFunc()
	}
	val, ok := p.vars[ident]
	if !ok {
		p.err = fmt.Errorf("variable '%s' is not bound", ident)
		val = 0
	}
	node := &variable{
		Key: ident,
		Val: val,
	}
	p.nextToken()
	return node
}

// 解析函数调用
func (p *parser) parseFunc() interfaces.INode {
	funcName := p.curToken.GetStr()
	p.nextToken() // skip LPAREN
	if p.curToken.GetType() != enums.LPAREN {
		p.err = fmt.Errorf("function '%s' expects '(' after name, got '%s'", funcName, p.curToken.GetStr())
		return nil
	}
	nodes := []interfaces.INode{}
	p.nextToken()
	for p.curToken.GetType() != enums.RPAREN && p.curToken.GetType() != enums.EOF {
		if p.curToken.GetType() == enums.COMMA {
			p.nextToken()
			continue
		}
		nodes = append(nodes, p.compile())
	}
	if p.curToken.GetType() != enums.RPAREN {
		p.err = fmt.Errorf("function '%s' expects ')' to close parameters, got '%s'", funcName, p.curToken.GetStr())
		return nil
	}
	fn, ok := p.fn[funcName]
	if !ok {
		p.err = fmt.Errorf("func '%s' is undefined", funcName)
		return nil
	}
	p.nextToken()
	return &function{
		Fn:   fn,
		Args: nodes,
	}
}

// 查看下一个 token
func (p *parser) peekNextToken() *lexer.Token {
	if p.index+1 >= len(p.tokens) {
		return lexer.EOF()
	}
	return &p.tokens[p.index+1]
}

// 获取下一个 token
func (p *parser) nextToken() *lexer.Token {
	p.index++
	if p.index >= len(p.tokens) {
		p.curToken = lexer.EOF()
	} else {
		p.curToken = &p.tokens[p.index]
	}
	return p.curToken
}

// 运算符优先级
func (p *parser) priority() int {
	switch p.curToken.GetType() {
	case enums.ADD, enums.SUB:
		return 1
	case enums.MUL, enums.QUO, enums.REM:
		return 2
	case enums.XOR:
		return 3
	default:
		return 0
	}
}
