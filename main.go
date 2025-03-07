package main

import (
	"fmt"

	"github.com/xupin/math-evaluate/lexer"
	"github.com/xupin/math-evaluate/parser"
)

func main() {
	lexer := lexer.New("(1+2)+(3+{FOUR}*4)")
	tokens, err := lexer.Lex()
	if err != nil {
		panic(err)
	}
	for _, token := range tokens {
		fmt.Printf("%+v \n", token)
	}
	p := parser.New(tokens)
	p.SetVar("FOUR", 4)
	ast, err := p.Parse()
	if err != nil {
		panic(err)
	}
	fmt.Printf("%.15f \n", ast.Evaluate())
}
