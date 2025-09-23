package mathevaluate_test

import (
	"testing"

	"github.com/xupin/math-evaluate/interfaces"
	"github.com/xupin/math-evaluate/lexer"
	"github.com/xupin/math-evaluate/parser"
)

func TestEvaluateExpression(t *testing.T) {
	input := "(1+2)+(3+FOUR*4)*random()"
	l := lexer.New(input)
	tokens, err := l.Lex()
	if err != nil {
		t.Fatalf("Lexer error: %v", err)
	}

	// 打印词法分析结果
	t.Log("Tokens:")
	for _, token := range tokens {
		t.Logf("%+v", token)
	}

	p := parser.New(tokens)
	p.SetVar("FOUR", 4)
	p.SetFunc("random", func(node ...interfaces.INode) float64 {
		return 1
	})

	ast, err := p.Parse()
	if err != nil {
		t.Fatalf("Parser error: %v", err)
	}

	result := ast.Evaluate()
	t.Logf("Evaluation result: %.15f", result)

	// 简单断言结果是否正确
	expected := 3 + (3 + 4*4) // (1+2) + (3+4*4) = 3 + 19 = 22
	if result != float64(expected) {
		t.Errorf("Expected %.15f, got %.15f", float64(expected), result)
	}
}
