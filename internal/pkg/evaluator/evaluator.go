package evaluator

import "github.com/elh/mal-go/internal/pkg/ast"

// evalAST looks up symbols in the given environment. This facilitates a mutual recursion between Eval and evalAST.
func evalAST(sexpr ast.Sexpr, env *Env) ast.Sexpr {
	switch sexpr.Type {
	case "list":
		var elems []ast.Sexpr
		for _, elem := range sexpr.Val.([]ast.Sexpr) {
			elems = append(elems, Eval(elem, env))
		}
		return ast.Sexpr{Type: "list", Val: elems}
	case "symbol":
		return env.Get(sexpr.Val.(string))
	default:
		return sexpr
	}
}

// Eval evaluates an s-expression in the given environment.
func Eval(expr ast.Sexpr, env *Env) ast.Sexpr {
	if expr.Type != "list" {
		return evalAST(expr, env)
	}
	list := expr.Val.([]ast.Sexpr)
	if len(list) == 0 {
		return expr
	}

	// def!
	if list[0].Type == "symbol" && list[0].Val.(string) == "def!" {
		if len(list) != 3 {
			panic("def! requires two arguments")
		}
		if list[1].Type != "symbol" {
			panic("def! requires a symbol as first argument")
		}
		v := Eval(list[2], env)
		env.Set(list[1].Val.(string), v)
		return v
	}

	// function call
	evaluatedList := evalAST(expr, env)
	elems := evaluatedList.Val.([]ast.Sexpr)
	if elems[0].Type != "function" {
		panic("first element of list must be a function")
	}

	fn := elems[0].Val.(func(args ...ast.Sexpr) ast.Sexpr)
	return fn(elems[1:]...)
}
