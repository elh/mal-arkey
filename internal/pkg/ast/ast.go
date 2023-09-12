package ast

// Sexpr is an s-expression with an associated type.
// Valid types:
// * "list"     - []ast.Sexpr
// * "symbol"   - string
// * "integer"  - int64
// * "function" - func(args ...ast.Sexpr) ast.Sexpr
type Sexpr struct {
	Type string
	Val  interface{}
}
