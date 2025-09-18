package parser

import (
	"log"
	"math"
	"math/rand/v2"
)

// 函数
var Funcs map[string]func(...node) float64

func InitFuncs() {
	Funcs = map[string]func(...node) float64{
		"min": func(n ...node) float64 {
			if len(n) < 2 {
				return 0
			}
			return math.Min(n[0].Evaluate(), n[1].Evaluate())
		},
		"max": func(n ...node) float64 {
			if len(n) < 2 {
				return 0
			}
			return math.Max(n[0].Evaluate(), n[1].Evaluate())
		},
		"floor": func(n ...node) float64 {
			if len(n) < 1 {
				return 0
			}
			return math.Floor(n[0].Evaluate())
		},
		"round": func(n ...node) float64 {
			if len(n) < 1 {
				return 0
			}
			return math.Round(n[0].Evaluate())
		},
		"random": func(n ...node) float64 {
			v := rand.Float64()
			log.Printf("random %f", v)
			return v
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
