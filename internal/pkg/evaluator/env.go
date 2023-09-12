package evaluator

import "github.com/elh/mal-go/internal/pkg/ast"

// Env is a map of symbols to bound values.
type Env struct {
	Outer *Env
	Data  map[string]ast.Sexpr
}

// GlobalEnv is the default global environment.
var GlobalEnv = Env{
	Outer: nil,
	Data: map[string]ast.Sexpr{
		"+": ast.Sexpr{Type: "function", Val: func(args ...ast.Sexpr) ast.Sexpr {
			var sum int64
			for _, arg := range args {
				sum += arg.Val.(int64)
			}
			return ast.Sexpr{Type: "integer", Val: sum}
		}},
		"-": ast.Sexpr{Type: "function", Val: func(args ...ast.Sexpr) ast.Sexpr {
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
		"*": ast.Sexpr{Type: "function", Val: func(args ...ast.Sexpr) ast.Sexpr {
			var product int64 = 1
			for _, arg := range args {
				product *= arg.Val.(int64)
			}
			return ast.Sexpr{Type: "integer", Val: product}
		}},
		"/": ast.Sexpr{Type: "function", Val: func(args ...ast.Sexpr) ast.Sexpr {
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
