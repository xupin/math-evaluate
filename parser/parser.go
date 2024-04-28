package parser

import (
	"errors"
	"fmt"
	"math"
	"strconv"

	"github.com/xupin/math-evaluate/enums"
	"github.com/xupin/math-evaluate/lexer"
)

type Parser struct {
	Tokens   []lexer.Token
	curToken *lexer.Token
	index    int
	err      error
	params   map[string]float64
}

type node interface {
	Evaluate() float64
}

type number struct {
	Val float64
}

type stmt struct {
	Type  int
	Left  node
	Right node
}

type variable struct {
	Key string
	Val float64
}

type function struct {
	Name string
	Args []node
}

// 函数
var Functions map[string]func(...node) float64

func init() {
	Functions = map[string]func(...node) float64{
		"min": func(n ...node) float64 {
			return math.Min(n[0].Evaluate(), n[1].Evaluate())
		},
		"max": func(n ...node) float64 {
			return math.Max(n[0].Evaluate(), n[1].Evaluate())
		},
		"floor": func(n ...node) float64 {
			return math.Floor(n[0].Evaluate())
		},
		"round": func(n ...node) float64 {
			return math.Round(n[0].Evaluate())
		},
	}
}

// 解析token
func (r *Parser) Parse() (node, error) {
	if len(r.Tokens) == 0 {
		return nil, errors.New("the token list is empty")
	}
	if r.curToken == nil {
		r.curToken = &r.Tokens[0]
	}
	return r.Compile(), r.err
}

// 设置变量
func (r *Parser) SetVar(key string, value float64) {
	if r.params == nil {
		r.params = make(map[string]float64, 0)
	}
	r.params[key] = value
}

// 构建树
func (r *Parser) Compile() node {
	left := r.ParseExpr()
	right := r.ParseRight(1, left)
	return right
}

// 从左开始处理
func (r *Parser) ParseExpr() node {
	switch r.curToken.GetType() {
	case enums.LPAREN:
		return r.ParseStmt()
	case enums.LBRACE:
		return r.ParseVar()
	case enums.NUMBER:
		return r.ParseNumber()
	case enums.ADD:
		return r.ParseNumber()
	case enums.SUB:
		if t := r.NextToken(); t.GetType() == enums.EOF {
			r.err = errors.New("expects to be number, eof given")
			return nil
		}
		return &stmt{
			Type:  enums.SUB,
			Left:  &number{},
			Right: r.ParseExpr(),
		}
	case enums.MUL:
		return r.ParseNumber()
	case enums.QUO:
		return r.ParseNumber()
	case enums.REM:
		return r.ParseNumber()
	case enums.XOR:
		return r.ParseNumber()
	case enums.FUNC:
		return r.ParseFunc()
	default:
		r.err = fmt.Errorf("expects to be number, '%s' given", r.curToken.GetStr())
		return nil
	}
}

// 处理操作符右侧
func (r *Parser) ParseRight(precedence int, left node) node {
	for {
		curPrec := r.Precedence()
		if curPrec < precedence {
			return left
		}
		tokenType := r.curToken.GetType()
		r.NextToken()
		right := r.ParseExpr()
		if right == nil {
			return nil
		}
		if curPrec < r.Precedence() {
			right = r.ParseRight(curPrec, right)
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
func (r *Parser) ParseVar() *variable {
	if t := r.NextToken(); t.GetType() == enums.EOF {
		r.err = errors.New("expects to be variable, eof given")
		return nil
	}
	key := r.curToken.GetStr()
	if t := r.NextToken(); t.GetType() == enums.EOF {
		r.err = errors.New("expects to be '}', eof given")
		return nil
	}
	v, ok := r.params[key]
	if !ok {
		r.err = fmt.Errorf("variable %s is not bound", key)
		v = 0
	}
	node := &variable{
		Key: key,
		Val: v,
	}
	r.NextToken()
	return node
}

// 数字
func (r *Parser) ParseNumber() *number {
	f, err := strconv.ParseFloat(r.curToken.GetStr(), 64)
	if err != nil {
		return &number{}
	}
	node := &number{
		Val: f,
	}
	r.NextToken()
	return node
}

// 表达式
func (r *Parser) ParseStmt() node {
	if t := r.NextToken(); t.GetType() == enums.EOF {
		r.err = errors.New("expects to be number, eof given")
		return nil
	}
	node := r.Compile()
	if node == nil {
		r.err = errors.New("expects to be number, nil given")
		return nil
	}
	if r.curToken.GetType() != enums.RPAREN {
		r.err = fmt.Errorf("expects to be number, '%s' given", r.curToken.GetStr())
		return nil
	}
	r.NextToken()
	return node
}

// 函数
func (r *Parser) ParseFunc() node {
	funcName := r.curToken.GetStr()
	if t := r.NextToken(); t.GetType() != enums.LPAREN {
		r.err = fmt.Errorf("expects to be '(', '%s' given", r.curToken.GetStr())
		return nil
	}
	nodes := make([]node, 0)
	for r.curToken.GetType() != enums.RPAREN {
		if t := r.NextToken(); t.GetType() == enums.EOF {
			break
		}
		if r.curToken.GetType() == enums.COMMA {
			continue
		}
		nodes = append(nodes, r.Compile())
	}
	if r.curToken.GetType() != enums.RPAREN {
		r.err = fmt.Errorf("expects to be number, '%s' given", r.curToken.GetStr())
		return nil
	}
	if _, ok := Functions[funcName]; !ok {
		r.err = fmt.Errorf("func %s is undefined", funcName)
		return nil
	}
	return &function{
		Name: funcName,
		Args: nodes,
	}
}

// 下一个token
func (r *Parser) NextToken() *lexer.Token {
	r.index++
	if r.index >= len(r.Tokens) {
		r.curToken = lexer.EOF()
	} else {
		r.curToken = &r.Tokens[r.index]
	}
	return r.curToken
}

// 优先级
func (r *Parser) Precedence() int {
	switch r.curToken.GetType() {
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

func (r *stmt) Evaluate() float64 {
	left := r.Left.Evaluate()
	right := r.Right.Evaluate()
	switch r.Type {
	case enums.ADD:
		return left + right
	case enums.SUB:
		return left - right
	case enums.MUL:
		return left * right
	case enums.QUO:
		if right == 0 {
			fmt.Printf("expr[%g/%g]exception, division by zero \n", left, right)
			return 0
		}
		return left / right
	case enums.REM:
		if right == 0 {
			fmt.Printf("expr[%g%%%g]exception, division by zero \n", left, right)
			return 0
		}
		return float64(int64(left) % int64(right))
	case enums.XOR:
		return math.Pow(left, right)
	default:
		return 0
	}
}

func (r *number) Evaluate() float64 {
	return r.Val
}

func (r *function) Evaluate() float64 {
	return Functions[r.Name](r.Args...)
}

func (r *variable) Evaluate() float64 {
	return r.Val
}

func (r *number) String() string {
	return fmt.Sprintf("{Type: %d, Val: %f}", enums.NUMBER, r.Val)
}

func (r *stmt) String() string {
	return fmt.Sprintf("{Type: %d, Left: %+v, Right: %+v}", r.Type, r.Left, r.Right)
}
