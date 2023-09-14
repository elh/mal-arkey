package mal

// Sexpr is an s-expression with an associated type.
// Valid types:
// * "list"         - []Sexpr
// * "symbol"       - string
// * "string"       - string
// * "integer"      - int64
// * "float"        - float64
// * "boolean"      - bool
// * "nil"          - nil
// * "atom"         - uuid atom table key. hack to dig myself out of non-pointer sexprs. i liked the immutability
// * "function"     - func(args ...Sexpr) Sexpr
// * "function-tco" - {
//   - AST:    Sexpr
//   - Params: Sexpr
//   - Env:    *Env
//   - Fn:     func(args ...Sexpr) Sexpr }
type Sexpr struct {
	Type string
	Val  interface{}
}

// FunctionTCO is a `fn*`-defined function that can be evaluated in a TCO style.
type FunctionTCO struct {
	AST    Sexpr
	Params []Sexpr
	Env    *Env
	Fn     func(args ...Sexpr) Sexpr
}
