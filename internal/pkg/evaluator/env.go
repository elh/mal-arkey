package evaluator

import "github.com/elh/mal-go/internal/pkg/ast"

// Env is a map of symbols to bound values.
type Env struct {
	Outer *Env
	Data  map[string]ast.Sexpr
}

// GlobalEnv creates a new default global environment.
func GlobalEnv() Env {
	return Env{
		Outer: nil,
		Data: map[string]ast.Sexpr{
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
		},
	}
}
