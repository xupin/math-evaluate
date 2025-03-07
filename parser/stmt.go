package parser

import (
	"fmt"
	"math"

	"github.com/xupin/math-evaluate/enums"
)

type stmt struct {
	Type  int
	Left  node
	Right node
}

func (s *stmt) String() string {
	return fmt.Sprintf("{Type: %d, Left: %+v, Right: %+v}", s.Type, s.Left, s.Right)
}

func (s *stmt) Evaluate() float64 {
	left := s.Left.Evaluate()
	right := s.Right.Evaluate()
	switch s.Type {
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
