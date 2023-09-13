package mal

// Sexpr is an s-expression with an associated type.
// Valid types:
// * "list"         - []ast.Sexpr
// * "symbol"       - string
// * "integer"      - int64
// * "float"        - float64
// * "boolean"      - bool
// * "nil"          - nil
// * "function"     - func(args ...ast.Sexpr) ast.Sexpr
// * "function-tco" - {
//   - ast:    ast.Sexpr
//   - params: ast.Sexpr
//   - env:    *evaluator.Env
//   - fn:     func(args ...ast.Sexpr) ast.Sexpr }
type Sexpr struct {
	Type string
	Val  interface{}
}
