package parser

type variable struct {
	Key string
	Val float64
}

func (v *variable) Evaluate() float64 {
	return v.Val
}
