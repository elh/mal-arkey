package mal

// TODO: change? I'm following the instructions from mal but I don't get this factoring.
func evalAST(sexpr Sexpr, env *Env) Sexpr {
	switch sexpr.Type {
	case "list":
		var elems []Sexpr
		for _, elem := range sexpr.Val.([]Sexpr) {
			elems = append(elems, Eval(elem, env))
		}
		return Sexpr{Type: "list", Val: elems}
	case "symbol":
		return env.Get(sexpr.Val.(string))
	default:
		return sexpr
	}
}

func evalDef(args []Sexpr, env *Env) Sexpr {
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

// With TCO. Return unevaluated body and new environment.
func evalLet(args []Sexpr, env *Env) (Sexpr, *Env) {
	letEnv := NewEnv(env, nil)
	if len(args) != 2 {
		panic("let* requires two arguments")
	}
	if args[0].Type != "list" {
		panic("let* requires a list as first argument")
	}
	bindings := args[0].Val.([]Sexpr)
	if len(bindings)%2 != 0 {
		panic("let* requires an even number of forms in bindings")
	}
	for i := 0; i < len(bindings); i += 2 {
		if bindings[i].Type != "symbol" {
			panic("let* bindings must be symbols")
		}
		letEnv.Set(bindings[i].Val.(string), Eval(bindings[i+1], letEnv))
	}

	return args[1], letEnv
}

// With TCO. Return unevaluated if/else branch form
func evalIf(args []Sexpr, env *Env) Sexpr {
	if len(args) != 2 && len(args) != 3 {
		panic("if requires three (or two) arguments")
	}
	cond := Eval(args[0], env)
	if (cond.Type == "boolean" && !cond.Val.(bool)) || (cond.Type == "nil") {
		if len(args) == 3 {
			return args[2]
		}
		return Sexpr{Type: "nil", Val: nil}
	}
	return args[1]
}

// With TCO. Return unevaluated final form
func evalDo(args []Sexpr, env *Env) Sexpr {
	for _, arg := range args[:len(args)-1] {
		Eval(arg, env)
	}
	return args[len(args)-1]
}

func evalFn(evalArgs []Sexpr, env *Env) Sexpr {
	if len(evalArgs) != 2 {
		panic("fn* requires two arguments")
	}
	if evalArgs[0].Type != "list" {
		panic("fn* requires a list as first argument")
	}
	params := evalArgs[0].Val.([]Sexpr)
	for _, param := range params {
		if param.Type != "symbol" {
			panic("fn* parameters must be symbols")
		}
	}
	body := evalArgs[1]

	return Sexpr{Type: "function", Val: func(args ...Sexpr) Sexpr {
		if len(params) != len(args) {
			panic("wrong number of arguments")
		}
		bindings := map[string]Sexpr{}
		for i, arg := range args {
			bindings[params[i].Val.(string)] = Eval(arg, env)
		}
		fnEnv := NewEnv(env, &bindings)
		return Eval(body, fnEnv)
	}}
}

// Eval evaluates an s-expression in the given environment.
func Eval(expr Sexpr, env *Env) Sexpr {
	// Tail call optimization prevents nested function calls.
	for {
		if expr.Type != "list" {
			return evalAST(expr, env)
		}
		list := expr.Val.([]Sexpr)
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
				expr, env = evalLet(args, env)
				continue
			case "if":
				expr = evalIf(args, env)
				continue
			case "do":
				expr = evalDo(args, env)
				continue
			case "fn*":
				return evalFn(args, env)
			}
		}

		// function call
		evaluatedList := evalAST(expr, env)
		elems := evaluatedList.Val.([]Sexpr)
		switch elems[0].Type {
		case "function":
			fn := elems[0].Val.(func(args ...Sexpr) Sexpr)
			return fn(elems[1:]...)
		case "function-tco":
			panic("TODO: unimplemented")
		default:
			panic("first element of list must be a function")
		}
	}
}
