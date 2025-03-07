package parser

import "math"

// 函数
var Funcs map[string]func(...node) float64

func InitFuncs() {
	Funcs = map[string]func(...node) float64{
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

type function struct {
	Name string
	Args []node
}

func (f *function) Evaluate() float64 {
	return Funcs[f.Name](f.Args...)
}
