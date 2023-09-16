package malarkey

// Value is a mal value with explicit type.
// Types:
// * "list"         - []Value
// * "vector"       - []Value
// * "hash-map"     - map[string]Value
// * "symbol"       - string
// * "string"       - string
// * "integer"      - int64
// * "float"        - float64
// * "boolean"      - bool
// * "nil"          - nil
// * "atom"         - int. atom id (`atoms` idx) hack to dig myself out of non-pointer vals. i liked the bias to immutability
// * "function"     - func(args ...Value) Value
// * "function-tco" - {
//   - AST:    Value
//   - Params: Value
//   - Env:    *Env
//   - Fn:     func(args ...Value) Value
//   - IsMacro bool                      }
type Value struct {
	Type string
	Val  interface{}
}

// FunctionTCO is a `fn*`-defined function that can be evaluated in a TCO style.
type FunctionTCO struct {
	AST     Value
	Params  []Value
	Env     *Env
	Fn      func(args ...Value) Value
	IsMacro bool
}
