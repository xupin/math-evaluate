package parser

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/xupin/math-evaluate/enums"
	"github.com/xupin/math-evaluate/lexer"
)

type parser struct {
	tokens   []lexer.Token
	curToken *lexer.Token
	index    int
	err      error
	params   map[string]float64
}

type node interface {
	Evaluate() float64
}

func init() {
	// 初始化函数
	InitFuncs()
}

func New(tokens []lexer.Token) *parser {
	return &parser{
		tokens: tokens,
	}
}

// 解析token
func (p *parser) Parse() (node, error) {
	if len(p.tokens) == 0 {
		return nil, errors.New("the token list is empty")
	}
	if p.curToken == nil {
		p.curToken = &p.tokens[0]
	}
	return p.compile(), p.err
}

// 设置变量
func (p *parser) SetVar(key string, value float64) {
	if p.params == nil {
		p.params = make(map[string]float64, 0)
	}
	p.params[key] = value
}

// 构建树
func (p *parser) compile() node {
	left := p.parseExpr()
	right := p.parseRight(1, left)
	return right
}

// 从左开始处理
func (p *parser) parseExpr() node {
	switch p.curToken.GetType() {
	case enums.LPAREN:
		return p.parseStmt()
	case enums.LBRACE:
		return p.parseVar()
	case enums.NUMBER:
		return p.parseNumber()
	case enums.ADD:
		return p.parseNumber()
	case enums.SUB:
		if t := p.nextToken(); t.GetType() == enums.EOF {
			p.err = errors.New("expects to be number, eof given")
			return nil
		}
		return &stmt{
			Type:  enums.SUB,
			Left:  &number{},
			Right: p.parseExpr(),
		}
	case enums.MUL:
		return p.parseNumber()
	case enums.QUO:
		return p.parseNumber()
	case enums.REM:
		return p.parseNumber()
	case enums.XOR:
		return p.parseNumber()
	case enums.FUNC:
		return p.parseFunc()
	default:
		p.err = fmt.Errorf("expects to be number, '%s' given", p.curToken.GetStr())
		return nil
	}
}

// 处理操作符右侧
func (p *parser) parseRight(priority int, left node) node {
	for {
		curPriority := p.priority()
		if curPriority < priority {
			return left
		}
		tokenType := p.curToken.GetType()
		p.nextToken()
		right := p.parseExpr()
		if right == nil {
			return nil
		}
		if curPriority < p.priority() {
			right = p.parseRight(curPriority, right)
			if right == nil {
				return nil
			}
		}
		left = &stmt{
			Type:  tokenType,
			Left:  left,
			Right: right,
		}
	}
}

// 变量
func (p *parser) parseVar() *variable {
	if t := p.nextToken(); t.GetType() == enums.EOF {
		p.err = errors.New("expects to be variable, eof given")
		return nil
	}
	key := p.curToken.GetStr()
	if t := p.nextToken(); t.GetType() == enums.EOF {
		p.err = errors.New("expects to be '}', eof given")
		return nil
	}
	v, ok := p.params[key]
	if !ok {
		p.err = fmt.Errorf("variable %s is not bound", key)
		v = 0
	}
	node := &variable{
		Key: key,
		Val: v,
	}
	p.nextToken()
	return node
}

// 数字
func (p *parser) parseNumber() *number {
	f, err := strconv.ParseFloat(p.curToken.GetStr(), 64)
	if err != nil {
		return &number{}
	}
	node := &number{
		Val: f,
	}
	p.nextToken()
	return node
}

// 表达式
func (p *parser) parseStmt() node {
	if t := p.nextToken(); t.GetType() == enums.EOF {
		p.err = errors.New("expects to be number, eof given")
		return nil
	}
	node := p.compile()
	if node == nil {
		p.err = errors.New("expects to be number, nil given")
		return nil
	}
	if p.curToken.GetType() != enums.RPAREN {
		p.err = fmt.Errorf("expects to be number, '%s' given", p.curToken.GetStr())
		return nil
	}
	p.nextToken()
	return node
}

// 函数
func (p *parser) parseFunc() node {
	funcName := p.curToken.GetStr()
	if t := p.nextToken(); t.GetType() != enums.LPAREN {
		p.err = fmt.Errorf("expects to be '(', '%s' given", p.curToken.GetStr())
		return nil
	}
	nodes := make([]node, 0)
	for p.curToken.GetType() != enums.RPAREN {
		if t := p.nextToken(); t.GetType() == enums.EOF {
			break
		}
		if p.curToken.GetType() == enums.COMMA {
			continue
		}
		nodes = append(nodes, p.compile())
	}
	if p.curToken.GetType() != enums.RPAREN {
		p.err = fmt.Errorf("expects to be number, '%s' given", p.curToken.GetStr())
		return nil
	}
	if _, ok := Funcs[funcName]; !ok {
		p.err = fmt.Errorf("func %s is undefined", funcName)
		return nil
	}
	return &function{
		Name: funcName,
		Args: nodes,
	}
}

// 下一个token
func (p *parser) nextToken() *lexer.Token {
	p.index++
	if p.index >= len(p.tokens) {
		p.curToken = lexer.EOF()
	} else {
		p.curToken = &p.tokens[p.index]
	}
	return p.curToken
}

// 优先级
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
