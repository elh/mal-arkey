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

func evalDef(args []ast.Sexpr, env *Env) ast.Sexpr {
	if len(args) != 2 {
		panic("def! requires two arguments")
	}
	if args[0].Type != "symbol" {
		panic("def! requires a symbol as first argument")
	}
	v := Eval(args[1], env)
	env.Set(args[0].Val.(string), v)
	return v
}

func evalLet(args []ast.Sexpr, env *Env) ast.Sexpr {
	letEnv := NewEnv(env, nil)
	if len(args) != 2 {
		panic("let* requires two arguments")
	}
	if args[0].Type != "list" {
		panic("let* requires a list as first argument")
	}
	bindings := args[0].Val.([]ast.Sexpr)
	if len(bindings)%2 != 0 {
		panic("let* requires an even number of forms in bindings")
	}
	for i := 0; i < len(bindings); i += 2 {
		if bindings[i].Type != "symbol" {
			panic("let* bindings must be symbols")
		}
		letEnv.Set(bindings[i].Val.(string), Eval(bindings[i+1], letEnv))
	}

	return Eval(args[1], letEnv)
}

func evalIf(args []ast.Sexpr, env *Env) ast.Sexpr {
	if len(args) != 2 && len(args) != 3 {
		panic("if requires three (or two) arguments")
	}
	cond := Eval(args[0], env)
	if (cond.Type == "boolean" && !cond.Val.(bool)) || (cond.Type == "nil") {
		if len(args) == 3 {
			return Eval(args[2], env)
		}
		return ast.Sexpr{Type: "nil", Val: nil}
	}
	return Eval(args[1], env)
}

func evalDo(args []ast.Sexpr, env *Env) ast.Sexpr {
	var result ast.Sexpr
	for _, arg := range args {
		result = Eval(arg, env)
	}
	return result
}

func evalFn(evalArgs []ast.Sexpr, env *Env) ast.Sexpr {
	if len(evalArgs) != 2 {
		panic("fn* requires two arguments")
	}
	if evalArgs[0].Type != "list" {
		panic("fn* requires a list as first argument")
	}
	params := evalArgs[0].Val.([]ast.Sexpr)
	for _, param := range params {
		if param.Type != "symbol" {
			panic("fn* parameters must be symbols")
		}
	}
	body := evalArgs[1]

	return ast.Sexpr{Type: "function", Val: func(args ...ast.Sexpr) ast.Sexpr {
		if len(params) != len(args) {
			panic("wrong number of arguments")
		}
		bindings := map[string]ast.Sexpr{}
		for i, arg := range args {
			bindings[params[i].Val.(string)] = Eval(arg, env)
		}
		fnEnv := NewEnv(env, &bindings)
		return Eval(body, fnEnv)
	}}
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

	// special forms
	if list[0].Type == "symbol" {
		args := list[1:]
		switch list[0].Val.(string) {
		case "def!":
			return evalDef(args, env)
		case "let*":
			return evalLet(args, env)
		case "if":
			return evalIf(args, env)
		case "do":
			return evalDo(args, env)
		case "fn*":
			return evalFn(args, env)
		}
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
