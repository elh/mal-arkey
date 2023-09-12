package evaluator

import (
	"fmt"

	"github.com/elh/mal-go/internal/pkg/ast"
	"github.com/elh/mal-go/internal/pkg/printer"
)

// Env is a map of symbols to bound values.
type Env struct {
	outer    *Env
	bindings map[string]ast.Sexpr
}

// NewEnv creates a new environment with the given outer environment.
func NewEnv(outer *Env, bindings *map[string]ast.Sexpr) *Env {
	if bindings == nil {
		bindings = &map[string]ast.Sexpr{}
	}
	return &Env{
		outer:    outer,
		bindings: *bindings,
	}
}

// Set binds a symbol to a value in the current environment.
func (e *Env) Set(symbol string, value ast.Sexpr) {
	e.bindings[symbol] = value
}

// Get returns the value bound to the given symbol in the environment.
func (e *Env) Get(symbol string) ast.Sexpr {
	if val, ok := e.bindings[symbol]; ok {
		return val
	}
	if e.outer != nil {
		return e.outer.Get(symbol)
	}
	panic(fmt.Sprintf("symbol '%v' not found", symbol))
}

func asFloat(v interface{}) float64 {
	switch v := v.(type) {
	case int64:
		return float64(v)
	case float64:
		return v
	default:
		panic(fmt.Sprintf("cannot convert %v to float64", v))
	}
}

// BuiltInEnv creates a new default built-in namespace env.
func BuiltInEnv() *Env {
	return &Env{
		outer: nil,
		bindings: map[string]ast.Sexpr{
			"+": {Type: "function", Val: func(args ...ast.Sexpr) ast.Sexpr {
				var sum int64
				for _, arg := range args {
					sum += arg.Val.(int64)
				}
				return ast.Sexpr{Type: "integer", Val: sum}
			}},
			"-": {Type: "function", Val: func(args ...ast.Sexpr) ast.Sexpr {
				if len(args) < 2 {
					panic("wrong number of arguments. `-` requires at least 2 arguments")
				}
				var diff int64
				for i, arg := range args {
					if i == 0 {
						diff = arg.Val.(int64)
					} else {
						diff -= arg.Val.(int64)
					}
				}
				return ast.Sexpr{Type: "integer", Val: diff}
			}},
			"*": {Type: "function", Val: func(args ...ast.Sexpr) ast.Sexpr {
				var product int64 = 1
				for _, arg := range args {
					product *= arg.Val.(int64)
				}
				return ast.Sexpr{Type: "integer", Val: product}
			}},
			"/": {Type: "function", Val: func(args ...ast.Sexpr) ast.Sexpr {
				if len(args) < 1 {
					panic("wrong number of arguments. `/` requires at least 1 arguments")
				}
				var quotient int64 = 1
				for i, arg := range args {
					if i == 0 && len(args) > 1 {
						quotient = arg.Val.(int64)
					} else {
						quotient /= arg.Val.(int64)
					}
				}
				return ast.Sexpr{Type: "integer", Val: quotient}
			}},
			"prn": {Type: "function", Val: func(args ...ast.Sexpr) ast.Sexpr {
				if len(args) < 1 {
					panic("wrong number of arguments. `prn` requires at least 1 arguments")
				}
				fmt.Println(printer.PrintStr(args[0]))
				return ast.Sexpr{Type: "nil", Val: nil}
			}},
			"list": {Type: "function", Val: func(args ...ast.Sexpr) ast.Sexpr {
				return ast.Sexpr{Type: "list", Val: args}
			}},
			"list?": {Type: "function", Val: func(args ...ast.Sexpr) ast.Sexpr {
				if len(args) > 0 && args[0].Type == "list" {
					return ast.Sexpr{Type: "boolean", Val: true}
				}
				return ast.Sexpr{Type: "boolean", Val: false}
			}},
			"empty?": {Type: "function", Val: func(args ...ast.Sexpr) ast.Sexpr {
				if len(args) > 0 && args[0].Type == "list" && len(args[0].Val.([]ast.Sexpr)) == 0 {
					return ast.Sexpr{Type: "boolean", Val: true}
				}
				return ast.Sexpr{Type: "boolean", Val: false}
			}},
			"count": {Type: "function", Val: func(args ...ast.Sexpr) ast.Sexpr {
				if len(args) > 0 && args[0].Type == "list" {
					return ast.Sexpr{Type: "integer", Val: int64(len(args[0].Val.([]ast.Sexpr)))}
				}
				return ast.Sexpr{Type: "integer", Val: int64(0)}
			}},
			"=": {Type: "function", Val: func(args ...ast.Sexpr) ast.Sexpr {
				if len(args) != 2 {
					panic("wrong number of arguments. `=` requires 2 arguments")
				}
				if args[0].Type == "list" && args[1].Type == "list" {
					alist := args[0].Val.([]ast.Sexpr)
					blist := args[1].Val.([]ast.Sexpr)
					if len(alist) != len(blist) {
						return ast.Sexpr{Type: "boolean", Val: false}
					}
					for i, a := range alist {
						if a.Type != blist[i].Type || a.Val != blist[i].Val {
							return ast.Sexpr{Type: "boolean", Val: false}
						}
					}
					return ast.Sexpr{Type: "boolean", Val: true}
				}

				if args[0].Type != args[1].Type || args[0].Val != args[1].Val {
					return ast.Sexpr{Type: "boolean", Val: false}
				}
				return ast.Sexpr{Type: "boolean", Val: true}
			}},
			"<": {Type: "function", Val: func(args ...ast.Sexpr) ast.Sexpr {
				if len(args) != 2 {
					panic("wrong number of arguments. `<` requires 2 arguments")
				}
				if asFloat(args[0].Val) < asFloat(args[1].Val) {
					return ast.Sexpr{Type: "boolean", Val: true}
				}
				return ast.Sexpr{Type: "boolean", Val: false}
			}},
			"<=": {Type: "function", Val: func(args ...ast.Sexpr) ast.Sexpr {
				if len(args) != 2 {
					panic("wrong number of arguments. `<` requires 2 arguments")
				}
				if asFloat(args[0].Val) <= asFloat(args[1].Val) {
					return ast.Sexpr{Type: "boolean", Val: true}
				}
				return ast.Sexpr{Type: "boolean", Val: false}
			}},
			">": {Type: "function", Val: func(args ...ast.Sexpr) ast.Sexpr {
				if len(args) != 2 {
					panic("wrong number of arguments. `<` requires 2 arguments")
				}
				if asFloat(args[0].Val) > asFloat(args[1].Val) {
					return ast.Sexpr{Type: "boolean", Val: true}
				}
				return ast.Sexpr{Type: "boolean", Val: false}
			}},
			">=": {Type: "function", Val: func(args ...ast.Sexpr) ast.Sexpr {
				if len(args) != 2 {
					panic("wrong number of arguments. `<` requires 2 arguments")
				}
				if asFloat(args[0].Val) >= asFloat(args[1].Val) {
					return ast.Sexpr{Type: "boolean", Val: true}
				}
				return ast.Sexpr{Type: "boolean", Val: false}
			}},
		},
	}
}
