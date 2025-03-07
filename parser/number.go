package parser

import (
	"fmt"

	"github.com/xupin/math-evaluate/enums"
)

type number struct {
	Val float64
}

func (n *number) String() string {
	return fmt.Sprintf("{Type: %d, Val: %f}", enums.NUMBER, n.Val)
}

func (n *number) Evaluate() float64 {
	return n.Val
}
