package ast

// Sexpr is an s-expression with an associated type.
type Sexpr struct {
	Type string
	Val  interface{}
}
