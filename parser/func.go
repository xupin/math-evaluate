package parser

import "github.com/xupin/math-evaluate/interfaces"

type function struct {
	Fn   funcType
	Args []interfaces.INode
}

func (f *function) Evaluate() float64 {
	return f.Fn(f.Args...)
}
